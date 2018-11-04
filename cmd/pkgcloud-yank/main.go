package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

	fmt.Printf("Yanking %s%d%s package(s) from %s ...\n", ansi.ColorCode("cyan"), len(packages), ansi.ColorCode("reset"), target)
	for _, pkg := range packages {
		go func(pkg string) {
			if err := client.Destroy(target.repo+"/"+target.distro, pkg); err != nil {
				errc <- fmt.Errorf("%s ... %s", pkg, err)
				return
			}
			fmt.Sprintf("%s/%s %s", target.repo, target.distro, pkg)
			resc <- fmt.Sprintf("%s ... OK", pkg)
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
