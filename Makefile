BUILDFLAGS=-a -ldflags "-w -s"

all: test build

test:
	go test -v $(go list ./... | grep -v /vendor/)

generate:
	go generate -x ./...

build:
	mkdir -p bin/
	go build -o bin/pkgcloud $(BUILDFLAGS) \
		github.com/tonylambiris/pkgcloud/cmd/pkgcloud

clean:
	rm -f bin/pkgcloud
