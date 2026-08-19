// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package fileutils_test

import (
	"fmt"
	"io/fs"
	"testing/fstest"

	"github.com/go-openapi/swag/fileutils"
)

// ExampleNewOverlayFS shows how to patch a set of assets shipped with a program,
// with a site-specific overlay that only holds what it changes.
//
// A file resolves in the topmost layer that holds it, so the overlay overrides the assets
// without having to repeat them. A directory reports the entries of every layer at once,
// so the templates of both layers are listed together.
func ExampleNewOverlayFS() {
	// the assets shipped with the program: in a real program, this could be an embed.FS
	shipped := fstest.MapFS{
		"config.yaml":          &fstest.MapFile{Data: []byte("theme: default")},
		"templates/index.html": &fstest.MapFile{Data: []byte("<h1>index</h1>")},
		"templates/about.html": &fstest.MapFile{Data: []byte("<h1>about</h1>")},
	}

	// a deployment overrides the configuration, and adds a template of its own
	site := fstest.MapFS{
		"config.yaml":            &fstest.MapFile{Data: []byte("theme: dark")},
		"templates/contact.html": &fstest.MapFile{Data: []byte("<h1>contact</h1>")},
	}

	assets := fileutils.NewOverlayFS(shipped, site)

	// the overlay wins for a file that both layers hold
	config, err := fs.ReadFile(assets, "config.yaml")
	if err != nil {
		fmt.Println("error:", err)

		return
	}
	fmt.Printf("config.yaml: %s\n", config)

	// a file that only the shipped assets hold still resolves
	about, err := fs.ReadFile(assets, "templates/about.html")
	if err != nil {
		fmt.Println("error:", err)

		return
	}
	fmt.Printf("about.html: %s\n", about)

	// the directory reports the templates of both layers, sorted by file name
	entries, err := fs.ReadDir(assets, "templates")
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("templates:")
	for _, entry := range entries {
		fmt.Println("  -", entry.Name())
	}

	// Output:
	// config.yaml: theme: dark
	// about.html: <h1>about</h1>
	// templates:
	//   - about.html
	//   - contact.html
	//   - index.html
}

// ExampleNewOpaqueFS shows how an overlay may take a directory over entirely,
// instead of merging its entries with the layers below.
//
// The deployment owns "templates", so the templates shipped with the program are neither
// listed nor readable. The rest of the tree keeps merging, so "config.yaml" still resolves.
func ExampleNewOpaqueFS() {
	shipped := fstest.MapFS{
		"config.yaml":          &fstest.MapFile{Data: []byte("theme: default")},
		"templates/index.html": &fstest.MapFile{Data: []byte("<h1>index</h1>")},
		"templates/about.html": &fstest.MapFile{Data: []byte("<h1>about</h1>")},
	}

	// this deployment provides the whole set of templates, and wants no other
	site := fstest.MapFS{
		"templates/index.html": &fstest.MapFile{Data: []byte("<h1>home</h1>")},
	}

	assets := fileutils.NewOverlayFS(shipped, fileutils.NewOpaqueFS(site, "templates"))

	entries, err := fs.ReadDir(assets, "templates")
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("templates:")
	for _, entry := range entries {
		fmt.Println("  -", entry.Name())
	}

	// the shipped template is hidden, not merely shadowed in the listing
	_, err = fs.ReadFile(assets, "templates/about.html")
	fmt.Println("about.html:", err)

	// the directories that are not owned keep merging
	config, err := fs.ReadFile(assets, "config.yaml")
	if err != nil {
		fmt.Println("error:", err)

		return
	}
	fmt.Printf("config.yaml: %s\n", config)

	// Output:
	// templates:
	//   - index.html
	// about.html: open templates/about.html: file does not exist
	// config.yaml: theme: default
}

// ExampleNewMapFS shows how to serve files held in memory, and stack them on top of another
// file system as an overlay.
//
// The map holds only what it overrides or adds: a directory of the base still reports the files
// the base holds beside them.
func ExampleNewMapFS() {
	shipped := fstest.MapFS{
		"templates/index.html": &fstest.MapFile{Data: []byte("<h1>index</h1>")},
		"templates/about.html": &fstest.MapFile{Data: []byte("<h1>about</h1>")},
	}

	// templates that a configuration provides, as raw bytes
	configured, err := fileutils.NewMapFS(fileutils.FromRawMap(map[string][]byte{
		"templates/index.html":   []byte("<h1>configured</h1>"),
		"templates/contact.html": []byte("<h1>contact</h1>"),
	}))
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	assets := fileutils.NewOverlayFS(shipped, configured)

	index, err := fs.ReadFile(assets, "templates/index.html")
	if err != nil {
		fmt.Println("error:", err)

		return
	}
	fmt.Printf("index.html: %s\n", index)

	entries, err := fs.ReadDir(assets, "templates")
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("templates:")
	for _, entry := range entries {
		fmt.Println("  -", entry.Name())
	}

	// Output:
	// index.html: <h1>configured</h1>
	// templates:
	//   - about.html
	//   - contact.html
	//   - index.html
}
