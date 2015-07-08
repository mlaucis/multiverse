// This package is here because godeps needs some help to get the dependencies right
package main

import (
	_ "github.com/tapglue/backend/config"
	_ "github.com/tapglue/backend/context"
	_ "github.com/tapglue/backend/server"
	_ "github.com/tapglue/backend/utils"
	_ "github.com/tapglue/backend/v02/server"
)

func main() {
	println("bye")
}
