package main

import (
	"log"
	"os"

	"github.com/sean9999/fasthak/templates"
)

func boilerplate() {

	sourceDir := templates.Vanilla()
	targetDir := "."

	err := os.CopyFS(targetDir, sourceDir)
	if err != nil {
		log.Fatal(err)
	}

}
