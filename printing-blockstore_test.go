package dex

import (
	"context"
	"testing"
	"time"

	blocks "github.com/ipfs/go-block-format"
	cid "github.com/ipfs/go-cid"
)

var testData = []byte(`Thy very songs not in thy songs,
No special strains to sing, none for itself,
But from the whole resulting, rising at last and floating,
A round full-orb'd eidolon.`)

var testDataB = []byte(`Green is the colour of her kind
Quickness of the eye
Deceives the mind
Envy is the bond between
The hopeful and the damned
`)

var testDataC = []byte(`We would sing and dance around
Because we know we can't be found
I'd like to be under the sea
In an octopus' garden in the shade
`)

// Helper function to check that the blockstore behaves as expected with
// and without preceding block puts
func helperDummyOps(t *testing.T, pbs *Pblockstore, testCid *cid.Cid, shouldHave bool) {
	has, err := pbs.Has(testCid)
	if has != shouldHave {
		t.Error("Has incorrectly reporting")
	}
	if err != nil {
		t.Error(err)
	}
	blks, err := pbs.Get(testCid)
	if err != ErrNotFound || blks != nil {
		t.Error("Get must always report not found")
	}
	err = pbs.DeleteBlock(testCid)
	if err != ErrNotFound {
		t.Error("DeleteBlock must always report not found")
	}
}

func newPbs() *Pblockstore {
	return &Pblockstore{
		membership: make(map[string]bool),
	}
}

// Test dummy (non-Put*) operations before Puts
func TestOpsCold(t *testing.T) {
	pbs := newPbs()
	// No panics from HashOnRead
	pbs.HashOnRead(true)
	pbs.HashOnRead(false)

	testBlock := blocks.NewBlock(testData)
	testCid := testBlock.Cid()
	helperDummyOps(t, pbs, testCid, false)
}

func TestAllKeysChan(t *testing.T) {
	pbs := newPbs()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	outchan, err := pbs.AllKeysChan(ctx)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-outchan:
		t.Error("should not be able to read from outchan")
	case <-ctx.Done():
		return
	}
}

func ExamplePut() {
	pbs := newPbs()
	testBlock := blocks.NewBlock(testData)
	pbs.Put(testBlock)
	// Output: [Block QmYmYZFATBaAWTGRL4Koe8hsHYFPwAKTYTqwWNH6Urp9sg]
}

func ExamplePutMany() {
	pbs := newPbs()
	testBlock := blocks.NewBlock(testData)
	testBlockB := blocks.NewBlock(testDataB)
	testBlockC := blocks.NewBlock(testDataC)
	pbs.PutMany([]blocks.Block{testBlock, testBlockB, testBlockC})
	// Output:
	// [Block QmYmYZFATBaAWTGRL4Koe8hsHYFPwAKTYTqwWNH6Urp9sg]
	// [Block QmdgQ9jTPNopx5fFbbYCRozkbMgsUm4zwBzGwn7tszPyAq]
	// [Block QmYkQssPzUkVxauJV7kJKPFdEqFQiGTuzZwtLSAWuRdrVE]
}

func TestOpsAfterPut(t *testing.T) {
	pbs := newPbs()
	testBlock := blocks.NewBlock(testData)
	testCid := testBlock.Cid()
	pbs.Put(testBlock)
	helperDummyOps(t, pbs, testCid, true)
}
