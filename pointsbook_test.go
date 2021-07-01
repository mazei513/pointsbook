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

func TestBookFromTransactions(t *testing.T) {
	b, err := pointsbook.BookFromTransactions("book-trx", []int{3, -1, 5})

	if err != nil {
		t.Fatalf("expected no error, got %s", err.Error())
	}
	assertBookCurrentPoints(t, b, 7)
	assertTrxLen(t, b, 3)
	assertUncommittedTrxLen(t, b, 0)
}

func TestBookAdd(t *testing.T) {
	b, err := pointsbook.NewBook("book-1")
	assertInitialBookState(t, b, err, "book-1")

	b.Add(1)

	assertBookCurrentPoints(t, b, 1)
	assertTrxLen(t, b, 1)
	assertUncommittedTrxLen(t, b, 1)
}

func BenchmarkCurrentPoints(b *testing.B) {
	bb, _ := pointsbook.NewBook("bench-book")
	for i := 0; i < 1000; i++ {
		bb.Add(uint(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bb.CurrentPoints()
	}
	b.StopTimer()
}

func TestBookSpendPoints(t *testing.T) {
	b, err := pointsbook.BookFromTransactions("book-b", []int{10})
	assertNoErr(t, err)
	assertBookCurrentPoints(t, b, 10)

	ok := b.Spend(5)
	if !ok {
		t.Fatalf("expected ok on spend, got not ok")
	}
	assertBookCurrentPoints(t, b, 5)
	assertTrxLen(t, b, 2)
	assertUncommittedTrxLen(t, b, 1)

	ok = b.Spend(6)
	if ok {
		t.Fatalf("expected not ok on spend, got ok")
	}
	assertBookCurrentPoints(t, b, 5)
	assertTrxLen(t, b, 2)
	assertUncommittedTrxLen(t, b, 1)
}

func assertInitialBookState(t *testing.T, b *pointsbook.Book, err error, id string) {
	t.Helper()
	assertNoErr(t, err)
	if b.ID() != id {
		t.Fatalf("expected id %s, got %s", strconv.Quote(id), strconv.Quote(b.ID()))
	}
	assertBookCurrentPoints(t, b, 0)
	assertTrxLen(t, b, 0)
	assertUncommittedTrxLen(t, b, 0)
}

func assertBookCurrentPoints(t *testing.T, b *pointsbook.Book, p int) {
	t.Helper()
	if b.CurrentPoints() != p {
		t.Fatalf("expected %d point(s), got %d", p, b.CurrentPoints())
	}
}
func assertTrxLen(t *testing.T, b *pointsbook.Book, n int) {
	t.Helper()
	if len(b.Transactions()) != n {
		t.Fatalf("expected %d transaction(s) got %d", n, len(b.Transactions()))
	}
}

func assertUncommittedTrxLen(t *testing.T, b *pointsbook.Book, n int) {
	t.Helper()
	if len(b.UncommittedTransactions()) != n {
		t.Fatalf("expected %d uncommitted transaction(s) got %d", n, len(b.UncommittedTransactions()))
	}
}

func assertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %s", err.Error())
	}
}
