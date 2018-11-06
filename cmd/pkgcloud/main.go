package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mgutz/ansi"
	"github.com/tonylambiris/pkgcloud"
)

var usage = "Usage: pkgcloud <push/yank> user/repo[/distro/version] /path/to/packages\n"

func main() {
	log.SetFlags(0)

	flag.Usage = func() { fmt.Fprintf(os.Stderr, usage) }
	flag.Parse()
	if flag.NArg() < 3 {
		log.Fatal(usage)
	}

	var f func(*pkgcloud.Client, *pkgcloud.Target, ...string) bool

	action := flag.Arg(0)

	switch action {
	case "push":
		f = actionPush
	case "yank":
		f = actionYank
	default:
		log.Fatalf("error: invalid action %s\n", action)
	}

	client, err := pkgcloud.NewClient("")
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	target, err := pkgcloud.NewTarget(flag.Arg(1))
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	packages := flag.Args()[2:]

	if f(client, target, packages...) {
		os.Exit(1)
	}
}

func actionPush(client *pkgcloud.Client, target *pkgcloud.Target, packages ...string) bool {
	resc := make(chan string)
	errc := make(chan error)

	fmt.Printf("Pushing %s%d%s package(s) to %s%s%s ...\n",
		ansi.Yellow, len(packages), ansi.Reset, ansi.Cyan, target, ansi.Reset,
	)

	for _, pkg := range packages {
		go func(pkg string) {
			pkgbase := path.Base(pkg)
			pkgname := fmt.Sprintf("%s%s%s", ansi.White, pkgbase, ansi.Reset)

			pkgs, err := client.Search(target.Repo, pkgbase, "", target.Distro, 0)
			if err != nil {
				errc <- fmt.Errorf("%s ... %s%s%s", pkgname, ansi.Magenta, err, ansi.Reset)
				return
			} else if len(pkgs) == 1 && pkgs[0].Filename == pkgbase {
				errc <- fmt.Errorf("%s ... %s%s%s", pkgname, ansi.Magenta, "package already exists", ansi.Reset)
				return
			}

			if err := client.CreatePackage(target.Repo, target.Distro, pkg); err != nil {
				errc <- fmt.Errorf("%s ... %s%s%s", pkgname, ansi.Magenta, err, ansi.Reset)
				return
			}
			resc <- fmt.Sprintf("%s ... %sOK%s", pkgname, ansi.Green, ansi.Reset)
		}(pkg)
	}

	failure := false
	for i, _ := range packages {
		select {
		case res := <-resc:
			log.Printf("%s(%d/%d)%s %s", ansi.Yellow, i+1, len(packages), ansi.Reset, res)
		case err := <-errc:
			log.Printf("%s(%d/%d)%s %s", ansi.Yellow, i+1, len(packages), ansi.Reset, err)
			failure = true
		}
	}

	return failure
}

func actionYank(client *pkgcloud.Client, target *pkgcloud.Target, packages ...string) bool {
	resc := make(chan string)
	errc := make(chan error)

	fmt.Printf("Yanking %s%d%s package(s) from %s%s%s ...\n",
		ansi.Yellow, len(packages), ansi.Reset, ansi.Cyan, target, ansi.Reset,
	)

	for _, pkg := range packages {
		go func(pkg string) {
			pkgbase := path.Base(pkg)
			pkgname := fmt.Sprintf("%s%s%s", ansi.White, pkgbase, ansi.Reset)

			if err := client.Destroy(target.String(), pkgbase); err != nil {
				errc <- fmt.Errorf("%s ... %s%s%s", pkgname, ansi.Magenta, err, ansi.Reset)
				return
			}
			resc <- fmt.Sprintf("%s ... %sOK%s", pkgname, ansi.Green, ansi.Reset)
		}(pkg)
	}

	failure := false
	for i, _ := range packages {
		select {
		case res := <-resc:
			log.Printf("%s(%d/%d)%s %s", ansi.Yellow, i+1, len(packages), ansi.Reset, res)
		case err := <-errc:
			log.Printf("%s(%d/%d)%s %s", ansi.Yellow, i+1, len(packages), ansi.Reset, err)
			failure = true
		}
	}

	return failure
}
