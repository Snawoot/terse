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
	if dst := r.AddViaIndex(); dst >= 0 {
		r.Load(dst, item)
	}
}

// Same as Add, but shifts actual write into datastructure to
// invoker in order to avoid copying element if it will not be added.
// Returns index where invoker should write to with Load method
// or -1 if element is not sampled. Invoker should put actual
// value of element only if non-negative value returned.
func (r *Reservoir[T]) AddViaIndex() int {
	newElem := chosenItem[T]{
		seqNo: r.consumed,
	}
	r.consumed++
	if newElem.seqNo < int64(r.size) {
		r.chosen = append(r.chosen, newElem)
		return len(r.chosen) - 1
	}
	potentialIdx := r.rng.Int63n(r.consumed + 1)
	if potentialIdx < int64(r.size) {
		r.chosen[int(potentialIdx)] = newElem
		return int(potentialIdx)
	}
	return -1
}

// Load sets actual value of item by index returned from AddViaIndex.
// That allows to postpone copy of added value and copy it only if element is
// actually sampled.
func (r *Reservoir[T]) Load(idx int, item T) {
	if idx < 0 {
		return
	}
	r.chosen[idx].item = item
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
