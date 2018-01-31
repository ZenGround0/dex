package dex

import (
	"context"
	"errors"
	"fmt"
	blocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
)

var ErrNotFound = errors.New("blockstore: block not found")

// Pblockstore is a dummy blockstore used for testing and sanity checks.
//
// It responds to Put and PutMany calls by printing out the blocks
// of data.
//
// Has returns false, even if Put was called on a block with this cid
//
// AllKeysChan blocks indefinitely
//
// Get returns an ErrNotFound even if block was previously Put
//
// DeleteBlock returns an ErrNotFound
//
// HashOnRead is a noop regardless of argument
type Pblockstore struct {
	membership map[string]bool
}

// Put prints out the block's representation string
func (pbs *Pblockstore) Put(block blocks.Block) error {
	id := block.Cid().String()
	pbs.membership[id] = true
	fmt.Printf("%s\n", block.String())
	return nil
}

// PutMany prints out each block
func (pbs *Pblockstore) PutMany(blocks []blocks.Block) error {
	for _, block := range blocks {
		err := pbs.Put(block)
		if err != nil {
			return err
		}
	}
	return nil
}

// Has returns false
func (pbs *Pblockstore) Has(c *cid.Cid) (bool, error) {
	_, ok := pbs.membership[c.String()]
	return ok, nil
}

// Get returns ErrNotFound
func (pbs *Pblockstore) Get(c *cid.Cid) (blocks.Block, error) {
	return nil, ErrNotFound
}

// DeleteBlock is a noop and returns ErrNotFound
func (pbs *Pblockstore) DeleteBlock(c *cid.Cid) error {
	return ErrNotFound
}

// HashOnRead is a noop
func (pbs *Pblockstore) HashOnRead(enabled bool) {
	return
}

// AllKeysChan blocks until the context expires
func (pbs *Pblockstore) AllKeysChan(ctx context.Context) (<-chan *cid.Cid, error) {
	output := make(chan *cid.Cid)
	go func() {
		defer func() {
			close(output)
		}()
		select {
		case <-ctx.Done():
			return
		}
	}()
	return output, nil
}
