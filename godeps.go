// This package is here because godeps needs some help to get the dependencies right
package main

import (
	_ "github.com/tapglue/multiverse/config"
	_ "github.com/tapglue/multiverse/context"
	_ "github.com/tapglue/multiverse/server"
	_ "github.com/tapglue/multiverse/utils"
	_ "github.com/tapglue/multiverse/v02/server"
	_ "github.com/tapglue/multiverse/v03/server"
	_ "github.com/tapglue/multiverse/v04/server"
)

func main() {
	println("bye")
}
