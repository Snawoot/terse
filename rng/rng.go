package rng

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"

	"pgregory.net/rand"
)

func NewRNG(seed *uint64) *rand.Rand {
	var initSeed uint64
	if seed != nil {
		initSeed = *seed
	} else {
		if err := binary.Read(crypto_rand.Reader, binary.BigEndian, &initSeed); err != nil {
			panic(fmt.Errorf("unable to get seed from CSPRNG: %w", err))
		}
	}
	return rand.New(initSeed)
}
