package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bayesianmind/demo-file-server/filestore"
	"github.com/bayesianmind/demo-file-server/server"
	"github.com/bayesianmind/demo-file-server/userstore"
)

func main() {
	fsPath := flag.String("fsPath", filepath.Join(os.TempDir(), "demo-fs"), "base path for user files")
	flag.Parse()
	fmt.Printf("Running file store on %q\n", *fsPath)
	serv := server.New(filestore.NewLocal(*fsPath), userstore.NewInMemory())
	err := serv.Run(":8081")
	if err != nil {
		panic(err)
	}
}
