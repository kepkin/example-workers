// +build generate

package main

import (
	"log"

	"github.com/kepkin/gorest"
)

func main() {
	err := gorest.Generate("openapi.yaml", gorest.Options{
		PackageName: "api",
		TargetFile:  "../api/api_gorest.go",
	})
	if err != nil {
		log.Fatal(err)
	}
}