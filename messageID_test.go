package slack

import (
	"testing"
)

func TestNewSafeID(t *testing.T) {
	idgen := NewSafeID(1)
	id1 := idgen.Next()
	id2 := idgen.Next()
	if id1 == id2 {
		t.Fatalf("id1 and id2 are same: id1: %d, id2: %d", id1, id2)
	}

	idgen = NewSafeID(100)
	id100 := idgen.Next()
	id101 := idgen.Next()
	if id2 == id100 {
		t.Fatalf("except id2 and id100 not same: id2: %d, id101: %d", id2, id100)
	}
	if id100 == id101 {
		t.Fatalf("id1 and id2 are same: id100: %d, id101: %d", id100, id101)
	}
}

var id int

func BenchmarkNewSafeID(b *testing.B) {
	b.ReportAllocs()

	idgen := NewSafeID(1)
	for i := 0; i < b.N; i++ {
		id = idgen.Next()
	}
}

func BenchmarkNewSafeIDParallel(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			idgen := NewSafeID(1)
			id = idgen.Next()
		}
	})
}
