package templates

import (
	"embed"
	"io/fs"
)

//go:embed vanilla
var vanilla embed.FS

func Vanilla() fs.FS {
	sub, _ := fs.Sub(vanilla, "vanilla")
	return sub
}
