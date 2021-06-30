package pointsbook_test

import (
	"strconv"
	"testing"

	"github.com/mazei513/pointsbook"
)

func TestNewBook(t *testing.T) {
	b := pointsbook.NewBook("book-id")

	assertInitialBookState(t, b, "book-id")
}

func TestBookAddPoints(t *testing.T) {
	b := pointsbook.NewBook("book-1")
	assertInitialBookState(t, b, "book-1")

	b.AddPoints(1)

	if b.CurrentPoints() != 1 {
		t.Fatalf("expected 1 point, got %d", b.CurrentPoints())
	}
	if len(b.Transactions()) != 1 {
		t.Fatalf("expected 1 transaction got %d", len(b.Transactions()))
	}
}

func assertInitialBookState(t *testing.T, b *pointsbook.Book, id string) {
	t.Helper()
	if b.ID() != id {
		t.Fatalf("expected id %s, got %s", strconv.Quote(id), strconv.Quote(b.ID()))
	}
	if b.CurrentPoints() != 0 {
		t.Fatalf("expected 0 points, got %d", b.CurrentPoints())
	}
	if len(b.Transactions()) != 0 {
		t.Fatalf("expected 0 transactions got %d", len(b.Transactions()))
	}
}

func BenchmarkCurrentPoints(b *testing.B) {
	bb := pointsbook.NewBook("bench-book")
	for i := 0; i < 1000; i++ {
		bb.AddPoints(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bb.CurrentPoints()
	}
	b.StopTimer()
}
