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

var usage = "Usage: pkgcloud-yank user/repo[/distro/version] /path/to/packages\n"

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

	fmt.Printf("Yanking %s%d%s package(s) from %s%s%s ...\n",
		ansi.Yellow, len(packages), ansi.Reset, ansi.Cyan, target, ansi.Reset,
	)
	for _, pkg := range packages {
		go func(pkg string) {
			remote := fmt.Sprintf("%s%s%s", ansi.White, path.Base(pkg), ansi.Reset)
			if err := client.Destroy(target.repo+"/"+target.distro, path.Base(pkg)); err != nil {
				errc <- fmt.Errorf("%s ... %s%s%s", remote, ansi.Magenta, err, ansi.Reset)
				return
			}
			resc <- fmt.Sprintf("%s ... %sOK%s", remote, ansi.Green, ansi.Reset)
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
