package pointsbook_test

import (
	"strconv"
	"testing"

	"github.com/mazei513/pointsbook"
)

func TestNewBook(t *testing.T) {
	t.Run("normal ID", func(t *testing.T) {
		b, err := pointsbook.NewBook("book-id")

		assertInitialBookState(t, b, err, "book-id")
	})
	t.Run("empty ID", func(t *testing.T) {
		b, err := pointsbook.NewBook("")

		if err == nil {
			t.Fatal("expected error, got none")
		}
		if b != nil {
			t.Fatal("expected nil Book, got non-nil")
		}
	})
}

func TestBookAddPoints(t *testing.T) {
	b, err := pointsbook.NewBook("book-1")
	assertInitialBookState(t, b, err, "book-1")

	b.AddPoints(1)

	if b.CurrentPoints() != 1 {
		t.Fatalf("expected 1 point, got %d", b.CurrentPoints())
	}
	if len(b.Transactions()) != 1 {
		t.Fatalf("expected 1 transaction got %d", len(b.Transactions()))
	}
}

func assertInitialBookState(t *testing.T, b *pointsbook.Book, err error, id string) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %s", err.Error())
	}
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
	bb, _ := pointsbook.NewBook("bench-book")
	for i := 0; i < 1000; i++ {
		bb.AddPoints(uint(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bb.CurrentPoints()
	}
	b.StopTimer()
}

func TestBookSpendPoints(t *testing.T) {
	b, err := pointsbook.NewBook("book-b")
	assertInitialBookState(t, b, err, "book-b")
	b.AddPoints(10)
	assertBookCurrentPoints(t, b, 10)

	ok := b.Spend(5)
	if !ok {
		t.Fatalf("expected ok on spend, got not ok")
	}
	assertBookCurrentPoints(t, b, 5)

	ok = b.Spend(6)
	if ok {
		t.Fatalf("expected not ok on spend, got ok")
	}
	assertBookCurrentPoints(t, b, 5)
}

func assertBookCurrentPoints(t *testing.T, b *pointsbook.Book, p int) {
	if b.CurrentPoints() != p {
		t.Fatalf("expected %d points, got %d", p, b.CurrentPoints())
	}
}
