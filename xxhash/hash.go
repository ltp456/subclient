package xxhash

import (
	"golang.org/x/crypto/blake2b"
	"hash"
)

type Blake2b128Concat struct {
	hasher hash.Hash
	data   []byte
}

// NewBlake2b128Concat returns an instance of blake2b concat hasher
func NewBlake2b128Concat(k []byte) (hash.Hash, error) {
	h, err := blake2b.New(16, k)
	if err != nil {
		return nil, err
	}
	return &Blake2b128Concat{hasher: h, data: k}, nil
}

// Write (via the embedded io.Writer interface) adds more data to the running hash.
func (bc *Blake2b128Concat) Write(p []byte) (n int, err error) {
	bc.data = append(bc.data, p...)
	return bc.hasher.Write(p)
}

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (bc *Blake2b128Concat) Sum(b []byte) []byte {
	return append(bc.hasher.Sum(b), bc.data...)
}

// Reset resets the Hash to its initial state.
func (bc *Blake2b128Concat) Reset() {
	bc.data = nil
	bc.hasher.Reset()
}

// Size returns the number of bytes Sum will return.
func (bc *Blake2b128Concat) Size() int {
	return len(bc.Sum(nil))
}

// BlockSize returns the hash's underlying block size.
// The Write method must be able to accept any amount
// of data, but it may operate more efficiently if all writes
// are a multiple of the block size.
func (bc *Blake2b128Concat) BlockSize() int {
	return bc.hasher.BlockSize()
}
