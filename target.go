package pkgcloud

import (
	"errors"
	"strings"
)

type Target struct {
	Repo   string
	Distro string
}

func NewTarget(s string) (*Target, error) {
	if s == "" {
		return nil, errors.New("empty target")
	}
	elems := strings.Split(s, "/")
	switch len(elems) {
	case 2:
		return &Target{
			Repo: strings.Join(elems[0:2], "/"),
		}, nil
	case 3:
		return &Target{
			Repo:   strings.Join(elems[0:2], "/"),
			Distro: elems[2],
		}, nil
	case 4:
		return &Target{
			Repo:   strings.Join(elems[0:2], "/"),
			Distro: strings.Join(elems[2:4], "/"),
		}, nil
	}
	return nil, errors.New("invalid target")
}

func (t Target) String() string {
	if t.Distro == "" {
		return t.Repo
	}
	return strings.Join([]string{t.Repo, t.Distro}, "/")
}
