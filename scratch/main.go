package main

import "log"

type typ interface {
	smth() string
}

func main() {
	t := *(new(typ))
	s := t.smth()
	log.Printf("%v", s)

	f := *getFn()
	err := f()
	log.Printf("%v", err)
}

type myFn func() error

func getFn() *myFn {
	return new(myFn)
}
