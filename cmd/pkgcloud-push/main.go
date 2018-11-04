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

var usage = "Usage: pkgcloud-push user/repo[/distro/version] /path/to/packages\n"

func main() {
	log.SetFlags(0)

	flag.Usage = func() { fmt.Fprintf(os.Stderr, usage) }
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal(usage)
	}

	target, err := newTarget(flag.Args()[0])
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	packages := flag.Args()[1:]

	client, err := pkgcloud.NewClient("")
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	resc := make(chan string)
	errc := make(chan error)

	fmt.Printf("Pushing %s%d%s package(s) to %s%s%s ...\n",
		ansi.Yellow, len(packages), ansi.Reset, ansi.Cyan, target, ansi.Reset,
	)
	for _, pkg := range packages {
		go func(pkg string) {
			pkgbase := path.Base(pkg)
			pkgname := fmt.Sprintf("%s%s%s", ansi.White, pkgbase, ansi.Reset)

			pkgs, err := client.Search(target.repo, pkgbase, "", target.distro, 0)
			if err != nil {
				errc <- fmt.Errorf("%s ... %s%s%s", pkgname, ansi.Magenta, err, ansi.Reset)
				return
			} else if len(pkgs) == 1 && pkgs[0].Filename == pkgbase {
				errc <- fmt.Errorf("%s ... %s%s%s", pkgname, ansi.Magenta, "package already exists", ansi.Reset)
				return
			}

			if err := client.CreatePackage(target.repo, target.distro, pkg); err != nil {
				errc <- fmt.Errorf("%s ... %s%s%s", pkgname, ansi.Magenta, err, ansi.Reset)
				return
			}
			resc <- fmt.Sprintf("%s ... %sOK%s", pkgname, ansi.Green, ansi.Reset)
		}(pkg)
	}

	failure := false
	for _, _ = range packages {
		select {
		case res := <-resc:
			log.Println(res)
		case err := <-errc:
			log.Println(err)
			failure = true
		}
	}
	if failure {
		os.Exit(1)
	}
}
