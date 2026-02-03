package panel

import (
	"embed"
	"io/fs"
)

//go:embed ui/*
var embedFS embed.FS

func GetFileSystem(useEmbed bool) (fs.FS, error) {
	if useEmbed {
		// Return the sub-filesystem anchored at "ui"
		return fs.Sub(embedFS, "ui")
	}
	return nil, nil // Caller should fallback to os.DirFS
}
