package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

func main() {
	httpFS, err := fsFROMOCI()
	if err != nil {
		panic(err)
	}
	fs := http.FileServer(httpFS)
	http.Handle("/", fs)
	log.Print("starting serving...")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		panic(err)
	}
}

func fsFROMOCI() (http.FileSystem, error) {
	ref, err := name.ParseReference("ghcr.io/frezbo/oci-static-website:0.0.1")
	if err != nil {
		return nil, err
	}
	img, err := remote.Image(ref)
	if err != nil {
		return nil, err
	}
	layers, err := img.Layers()
	if err != nil {
		return nil, err
	}
	uncompressed, err := layers[0].Uncompressed()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(uncompressed)
	if err != nil {
		return nil, err
	}
	defer uncompressed.Close()
	vfs := mapfs.New(map[string]string{"index.html": buf.String()})
	httpFS := httpfs.New(vfs)
	return httpFS, nil
}
