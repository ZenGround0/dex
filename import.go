package dex

import (
	"context"
	"strings"

	"github.com/ipfs/go-ipfs-cmdkit/files"
	"github.com/ipfs/go-ipfs/core/coreunix"
)

// ImportToPrint is a first step towards a streaming importer.  Verify that we
// can hijack the blockstore abstraction to redirect blocks as they arrive
// Closely follows go-ipfs/core/commands/add.go: Run func
func ImportToPrint(file files.File) error {
	dserv := &pDAGService{
		membership: make(map[string]bool),
	}

	ctx := context.Background() // using background for now, should upgrade later
	fileAdder, err := coreunix.NewAdder(ctx, nil, nil, dserv)
	if err != nil {
		return err
	}
	fileAdder.Pin = false // This way we can be honest that blockstore doesn't exist

	// add the file
	if err := fileAdder.AddFile(file); err != nil {
		return err
	}

	// Without this call all of the directory nodes (stored in MFS) do not get
	// written through to the dagservice and its blockstore

	_, err = fileAdder.Finalize()
	// ignore errors caused by printing-blockstore Get not finding blocks
	if !strings.Contains(err.Error(), "dagservice: block not found") {
		return err
	}
	return nil

	// Output is exfiltrated from within the blockstore (here it is printed)
	// when we have a streaming blockstore we will want to include a channel
	// as an arg to this function and a param to blockstore init so that output
	// channel can be registered
	//
	// TODO: Will the streaming blockstore use a channel or something else, like
	// libp2p streams?
}
