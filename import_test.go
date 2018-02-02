package dex

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/ipfs/go-ipfs-cmdkit/files"
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
