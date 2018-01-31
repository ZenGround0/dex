package dex

import (
	"context"

	"github.com/ipfs/go-ipfs-cmdkit/files"
	bstore "github.com/ipfs/go-ipfs/blocks/blockstore"
	blockservice "github.com/ipfs/go-ipfs/blockservice"
	"github.com/ipfs/go-ipfs/core/coreunix"
	"github.com/ipfs/go-ipfs/exchange/offline"
	dag "github.com/ipfs/go-ipfs/merkledag"
)

// A first step to a streaming importer.  Verify that we can hijack the
// blockstore abstraction to redirect blocks as they arrive
// Closely follows go-ipfs/core/commands/add.go: Run func
func ImportToPrint(file files.File) error {
	// Init objects needed by adder
	// DAGSERVICE [√]
	// blockservice [√]
	// GC-BLOCKSTORE[√]
	// printing blockstore [√]
	// dummy GC locker [√] -- normal GCLocker
	// dummy exchange [√]   -- offline.exchange
	// DUMMY PINNING [√]  -- nil for now
	pbs := &Pblockstore{
		membership: make(map[string]bool),
	} // This "stores" blocks by printing them to stdout
	locker := bstore.NewGCLocker()
	addblockstore := bstore.NewGCBlockstore(pbs, locker)

	exch := offline.Exchange(addblockstore)
	bserv := blockservice.New(addblockstore, exch)
	dserv := dag.NewDAGService(bserv)
	// TODO: confirm GC should not ever be called on these runs, or come up with
	// a pinner that works in tandem with the printint blockstore.
	//   -- I think this is safe as no-one else has a blockstore ref to call GCLock
	//   -- one way to ensure this is prevented is by constructing our own
	//      dummy GCLocker that always returns false on GCRequested()
	// pinning := nil
	ctx := context.Background() // using background for now, should upgrade later
	fileAdder, err := coreunix.NewAdder(ctx, nil, addblockstore, dserv)
	if err != nil {
		return err
	}

	// add the file
	if err := fileAdder.AddFile(file); err != nil {
		return err
	}

	// Without this call all of the directory nodes (stored in MFS) do not get
	// written through to the dagservice and its blockstore
	_, err = fileAdder.Finalize()
	return err

	// Output is exfiltrated from within the blockstore (here it is printed)
	// when we have a streaming blockstore we will want to include a channel
	// as an arg to this function and a param to blockstore init so that output
	// channel can be registered
	//
	// TODO: Will the streaming blockstore use a channel or something else, like
	// libp2p streams?
}
