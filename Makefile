all: test generate build

test:
	go test -v ./...

generate:
	go generate -x ./...

build:
	mkdir -p bin/
	go build -o bin/pkgcloud-push \
		-ldflags "-w -s" github.com/tonylambiris/pkgcloud/cmd/pkgcloud-push
	go build -o bin/pkgcloud-yank \
		-ldflags "-w -s" github.com/tonylambiris/pkgcloud/cmd/pkgcloud-yank

clean:
	rm -f bin/pkgcloud-push bin/pkgcloud-yank
