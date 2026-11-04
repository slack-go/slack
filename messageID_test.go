package slack

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSafeID(t *testing.T) {
	for _, startID := range []int{-100, 0, 1, 100} {
		idgen := NewSafeID(startID)
		for offset := range 2 {
			want := startID + offset + 1
			require.Equal(t, want, idgen.Next(), "NewSafeID(%d)", startID)
		}
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
