package aerospike

import "math/rand"

// RandomToken returns a random Token
func RandomID() int64 {
	return rand.Int63()
}
