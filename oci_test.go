package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func TestOCI(t *testing.T) {
	ref, err := name.ParseReference("ghcr.io/frezbo/oci-static-website:0.0.1")
	if err != nil {
		t.Error(err)
	}
	img, err := remote.Image(ref)
	if err != nil {
		t.Error(err)
	}
	layers, err := img.Layers()
	if err != nil {
		t.Error(err)
	}
	buf := new(bytes.Buffer)

	compressed, err := layers[0].Uncompressed()
	if err != nil {
		t.Error(err)
	}

	_, err = buf.ReadFrom(compressed)
	if err != nil {
		t.Error(err)
	}
	defer compressed.Close()
	fmt.Println(buf.String())
}
