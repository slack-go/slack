package slack

import "sync/atomic"

// IDGenerator provides an interface for generating integer ID values.
type IDGenerator interface {
	Next() int
}

// NewSafeID returns a new instance of an IDGenerator which is safe for
// concurrent use by multiple goroutines.
func NewSafeID(startID int) IDGenerator {
	return &safeID{
		nextID: int64(startID),
	}
}

type safeID struct {
	nextID int64
}

func (s *safeID) Next() (id int) {
	id = int(atomic.LoadInt64(&s.nextID))
	atomic.AddInt64(&s.nextID, 1)

	return id
}
