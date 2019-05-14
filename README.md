# pkgcloud

Talk to the [packagecloud API](https://packagecloud.io/docs/api), in Go.

## Installation

    $ go get github.com/penguinpowernz/pkgcloud/...

## API Usage

See [Godoc](https://godoc.org/github.com/penguinpowernz/pkgcloud) and
[pkgcloud.go](pkgcloud.go) to learn about the API.

## Client Usage

### Pushing/Yanking packages

This tool is a simple and fast replacement for the original `package_cloud`
command. If you pass more than one package, `pkgcloud push` will push them in
parallel! Before using it, however, make sure that `PACKAGECLOUD_TOKEN` is set
in your environment:

    alias pkgcloud-push='PACKAGECLOUD_TOKEN=0xDEADBEEF pkgcloud push myaccount/myrepo/el/7'
    alias pkgcloud-yank='PACKAGECLOUD_TOKEN=0xDEADBEEF pkgcloud yank myaccount/myrepo/el/7'

Usage:

    $ pkgcloud <push/yank> user/repo[/distro/version] /path/to/packages

Examples:

    # Debian
    $ pkgcloud push myaccount/myrepo/ubuntu/trusty example_1.2.3_amd64.deb

    # RPM
    $ pkgcloud push myaccount/myrepo/el/7 *.rpm
    $ pkgcloud yank myaccount/myrepo/el/7 *.rpm

    # RubyGem
    $ pkgcloud push myaccount/myrepo example-1.2.3.gem
