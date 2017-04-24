package main

import (
	"flag"

	"github.com/c00w/demesure/lib"
)

var listen = flag.Bool("listen", false, "Should we listen?")

func main() {

	flag.Parse()

	lib.DoIt(*listen, flag.Arg(0))

}
