package reservoir

import "sort"

// RNG compatible with math/rand PRNG.
type RNG interface {
	Int63n(n int64) int64
}

// Implements reservoir sampling algorithm, but also retains 
// original order of elements.
type Reservoir[T any] struct {
	size     int
	rng      RNG
	consumed int64
	chosen   []chosenItem[T]
}

type chosenItem[T any] struct {
	seqNo int64
	item  T
}

// Created new Reservoir instance. Function will panic if size is negative.
func NewReservoir[T any](size int, rng RNG) *Reservoir[T] {
	if size < 0 {
		panic("negative reservoir size")
	}
	return &Reservoir[T]{
		size:   size,
		rng:    rng,
		chosen: make([]chosenItem[T], 0, size),
	}
}

// Adds another candidate item to reservoir
func (r *Reservoir[T]) Add(item T) {
	newElem := chosenItem[T]{
		seqNo: r.consumed,
		item:  item,
	}
	r.consumed++
	if newElem.seqNo < int64(r.size) {
		r.chosen = append(r.chosen, newElem)
		return
	}
	potentialIdx := r.rng.Int63n(r.consumed + 1)
	if potentialIdx < int64(r.size) {
		r.chosen[int(potentialIdx)] = newElem
	}
}

// Gets sampled items
func (r *Reservoir[T]) Items() []T {
	chosen := make([]chosenItem[T], len(r.chosen))
	copy(chosen, r.chosen)
	sort.Slice(chosen, func(i, j int) bool {
		return chosen[i].seqNo < chosen[j].seqNo
	})
	result := make([]T, len(chosen))
	for i := range chosen {
		result[i] = chosen[i].item
	}
	return result
}
