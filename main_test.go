package main

import (
	"os"
	"testing"

	packr "github.com/gobuffalo/packr/v2"
	"github.com/gregdhill/go-openrpc/generate"
	"github.com/gregdhill/go-openrpc/parse"
)

func generateExampleProxyServer() error {
	specFile := "parse/testdata/eth_openrpc.json"
	pkgDir := "rpc"

	openrpc, err := readSpec(specFile)
	if err != nil {
		return err
	}

	parse.GetTypes(openrpc, openrpc.Objects)
	box := packr.New("template", "./templates")

	if err = generate.WriteFile(box, "server", pkgDir, openrpc); err != nil {
		return err
	}

	if err = generate.WriteFile(box, "types", pkgDir, openrpc); err != nil {
		return err
	}
	if err = generate.WriteFile(box, "example-proxy-server", "main", openrpc); err != nil {
		return err
	} else {
		// HACK
		if err := os.MkdirAll("example/", os.ModePerm); err != nil {
			return err
		}
		if err := os.Rename("main/example-proxy-server.go", "example/example-proxy-server.go"); err != nil {
			return err
		}
		if err := os.RemoveAll("main/"); err != nil {
			return err
		}
	}
	return nil
}

func TestExampleProxyServer(t *testing.T) {
	err := generateExampleProxyServer()
	if err != nil {
		t.Fatal(err)
	}
}
