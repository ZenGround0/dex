package dex

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/ipfs/go-ipfs-cmdkit/files"
	ipld "github.com/ipfs/go-ipld-format"
)

const testDir = "testingData"

func getTestingDir() (files.File, error) {
	fpath := testDir
	stat, err := os.Lstat(fpath)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("testDir should be seen as directory")
	}

	return files.NewSerialFile(path.Base(fpath), fpath, false, stat)
}

// simply import and ensure no errors occur
func TestImportPrint(t *testing.T) {
	file, err := getTestingDir()
	if err != nil {
		t.Fatal(err)
	}
	err = ImportToPrint(file)
	if err != nil {
		t.Fatal(err)
	}
}

// import and receive all blocks
func TestImportChannel(t *testing.T) {
	file, err := getTestingDir()
	if err != nil {
		t.Fatal(err)
	}

	outChan := make(chan *ipld.Node)
	go func() {
		for node := range outChan {
			fmt.Printf("%s\n", (*node).String())
		}
	}()
	err = ImportToChannel(file, outChan, context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
