# fanal
Static Analysis Library for Containers

[![GoDoc](https://godoc.org/github.com/BryanKMorrow/fanal?status.svg)](https://godoc.org/github.com/BryanKMorrow/fanal)
![Test](https://github.com/BryanKMorrow/fanal/workflows/Test/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/BryanKMorrow/fanal)](https://goreportcard.com/report/github.com/BryanKMorrow/fanal)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/BryanKMorrow/fanal/blob/main/LICENSE)

## Feature
- Detect OS
- Extract OS packages
- Extract libraries used by an application
  - Bundler, Composer, npm, Yarn, Pipenv, Poetry, Cargo

## Example
See [`cmd/fanal/`](cmd/fanal)

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/xerrors"

	"github.com/BryanKMorrow/fanal/cache"

	"github.com/BryanKMorrow/fanal/analyzer"
	_ "github.com/BryanKMorrow/fanal/analyzer/library/bundler"
	_ "github.com/BryanKMorrow/fanal/analyzer/library/composer"
	_ "github.com/BryanKMorrow/fanal/analyzer/library/npm"
	_ "github.com/BryanKMorrow/fanal/analyzer/library/pipenv"
	_ "github.com/BryanKMorrow/fanal/analyzer/library/poetry"
	_ "github.com/BryanKMorrow/fanal/analyzer/library/yarn"
	_ "github.com/BryanKMorrow/fanal/analyzer/library/cargo"
	_ "github.com/BryanKMorrow/fanal/analyzer/os/alpine"
	_ "github.com/BryanKMorrow/fanal/analyzer/os/amazonlinux"
	_ "github.com/BryanKMorrow/fanal/analyzer/os/debianbase"
	_ "github.com/BryanKMorrow/fanal/analyzer/os/suse"
	_ "github.com/BryanKMorrow/fanal/analyzer/os/redhatbase"
	_ "github.com/BryanKMorrow/fanal/analyzer/pkg/apk"
	_ "github.com/BryanKMorrow/fanal/analyzer/pkg/dpkg"
	_ "github.com/BryanKMorrow/fanal/analyzer/pkg/rpm"
	"github.com/BryanKMorrow/fanal/extractor"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	ctx := context.Background()
	tarPath := flag.String("f", "-", "layer.tar path")
	clearCache := flag.Bool("clear", false, "clear cache")
	flag.Parse()

	if *clearCache {
		if err = cache.Clear(); err != nil {
			return xerrors.Errorf("error in cache clear: %w", err)
		}
	}

	args := flag.Args()

	var files extractor.FileMap
	if len(args) > 0 {
		files, err = analyzer.Analyze(ctx, args[0])
		if err != nil {
			return err
		}
	} else {
		rc, err := openStream(*tarPath)
		if err != nil {
			return err
		}

		files, err = analyzer.AnalyzeFromFile(ctx, rc)
		if err != nil {
			return err
		}
	}

	os, err := analyzer.GetOS(files)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", os)

	pkgs, err := analyzer.GetPackages(files)
	if err != nil {
		return err
	}
	fmt.Printf("Packages: %d\n", len(pkgs))

	libs, err := analyzer.GetLibraries(files)
	if err != nil {
		return err
	}
	for filepath, libList := range libs {
		fmt.Printf("%s: %d\n", filepath, len(libList))
	}
	return nil
}

func openStream(path string) (*os.File, error) {
	if path == "-" {
		if terminal.IsTerminal(0) {
			flag.Usage()
			os.Exit(64)
		} else {
			return os.Stdin, nil
		}
	}
	return os.Open(path)
}

```


## Notes
When using `latest` tag, that image will be cached. After `latest` tag is updated, you need to clear cache.



