package cache

import (
	"fmt"
)

type Key interface {
	String() string
}

type keyWrapper struct {
	ns  string
	key Key
}

func newKeyWrapper(ns string, key Key) keyWrapper {
	return keyWrapper{
		ns:  ns,
		key: key,
	}
}

func (k keyWrapper) String() string {
	return fmt.Sprintf("%s:%s", k.ns, k.key.String())
}
