//go:build ignore

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

// this is invoked from the server pkg level
func main() {
	if err := entc.Generate("./internal/ent/schema", &gen.Config{}); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
