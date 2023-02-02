package reservoir

import "sort"

type RNG interface {
	Int63n(n int64) int64
}

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
