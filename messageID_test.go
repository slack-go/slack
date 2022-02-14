package slack

import (
	"testing"
)

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
