package dex

import (
	"context"

	blocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
)

// nop-blockstore is a dummy blockstore used to satisfy interfaces
// without interfering with import operations
//
// Operations are errorless nops
type nopBlockstore struct{}

// Put does nothing
func (nop *nopBlockstore) Put(block blocks.Block) error {
	return nil
}

// PutMany does nothing
func (nop *nopBlockstore) PutMany(blocks []blocks.Block) error {
	return nil
}

// Has returns false and no error
func (nop *nopBlockstore) Has(c *cid.Cid) (bool, error) {
	return false, nil
}

// Get returns nil and no error
func (nop *nopBlockstore) Get(c *cid.Cid) (blocks.Block, error) {
	return nil, nil
}

// DeleteBlock is a noop and returns nil
func (nop *nopBlockstore) DeleteBlock(c *cid.Cid) error {
	return nil
}

// HashOnRead does nothing
func (nop *nopBlockstore) HashOnRead(enabled bool) {
	return
}

// AllKeysChan blocks until the context expires
func (nop *nopBlockstore) AllKeysChan(ctx context.Context) (<-chan *cid.Cid, error) {
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
