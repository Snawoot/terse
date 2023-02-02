package reservoir

import (
	"math/rand"
	"testing"
)

func TestConstructor(t *testing.T) {
	_ = NewReservoir[string](1, rand.New(rand.NewSource(0)))
}

func TestIncomplete1(t *testing.T) {
	r := NewReservoir[int](10, rand.New(rand.NewSource(0)))
	for i := 0; i < 10; i++ {
		r.Add(i)
	}
	res := r.Items()
	for i := 0; i < 10; i++ {
		if res[i] != i {
			t.Fatalf("unexpected output: res[%d] is %d, but %d is expected. res = %#v", i, res[i], i, res)
		}
	}
}

func TestIncomplete2(t *testing.T) {
	r := NewReservoir[int](10, rand.New(rand.NewSource(0)))
	for i := 0; i < 3; i++ {
		r.Add(i)
	}
	res := r.Items()
	for i := 0; i < 3; i++ {
		if res[i] != i {
			t.Fatalf("unexpected output: res[%d] is %d, but %d is expected. res = %#v", i, res[i], i, res)
		}
	}
}

func TestOverflow(t *testing.T) {
	r := NewReservoir[int](10, rand.New(rand.NewSource(0)))
	for i := 0; i < 11; i++ {
		r.Add(i)
	}
	res := r.Items()
	if len(res) != 10 {
		t.Fatalf("unexpected lenght: %d", len(res))
	}
}
