package rng

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	math_rand "math/rand"
)

func NewRNG(seed *int64) *math_rand.Rand {
	var initSeed int64
	if seed != nil {
		initSeed = *seed
	} else {
		if err := binary.Read(crypto_rand.Reader, binary.BigEndian, &initSeed); err != nil {
			panic(fmt.Errorf("unable to get seed from CSPRNG: %w", err))
		}
	}
	return math_rand.New(math_rand.NewSource(initSeed))
}
