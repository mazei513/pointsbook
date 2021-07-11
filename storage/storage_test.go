package storage_test

import (
	"context"
	"errors"
	"path"
	"strconv"
	"testing"

	"github.com/mazei513/pointsbook"
	"github.com/mazei513/pointsbook/storage"
)

func TestNewStore(t *testing.T) {
	dir := t.TempDir()
	s, err := storage.NewStore(path.Join(dir, "TestNewStore.db"))

	assertNoErr(t, err)
	if s == nil {
		t.Fatal("expected non-nil store, got nil")
	}

	s.Close()
}

func TestMigrateStore(t *testing.T) {
	dir := t.TempDir()
	s, err := storage.NewStore(path.Join(dir, "TestMigrateStore.db"))
	assertNoErr(t, err)
	defer s.Close()
	ctx := context.Background()

	t.Run("uninitialized", func(t *testing.T) {
		v, err := s.SchemaVersion(ctx)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, storage.ErrUninitialized) {
			t.Fatalf("expected %s error, got %s", strconv.Quote(storage.ErrUninitialized.Error()), strconv.Quote(err.Error()))
		}
		if v != 0 {
			t.Fatalf("expected version 0, got %d", v)
		}
	})

	t.Run("uninitialized to 0", func(t *testing.T) {
		err := s.MigrateTo(ctx, 0)
		assertNoErr(t, err)

		v, err := s.SchemaVersion(ctx)
		assertNoErr(t, err)
		if v != 0 {
			t.Fatalf("expected version 0, got %d", v)
		}
	})

	t.Run("same version", func(t *testing.T) {
		err := s.MigrateTo(ctx, 0)
		assertNoErr(t, err)

		v, err := s.SchemaVersion(ctx)
		assertNoErr(t, err)
		if v != 0 {
			t.Fatalf("expected version 0, got %d", v)
		}
	})

	t.Run("backwards", func(t *testing.T) {
		err := s.MigrateTo(ctx, -1)
		assertErr(t, err)

		v, err := s.SchemaVersion(ctx)
		assertNoErr(t, err)
		if v != 0 {
			t.Fatalf("expected version 0, got %d", v)
		}
	})

	t.Run("0 to 1", func(t *testing.T) {
		err := s.MigrateTo(ctx, 1)
		assertNoErr(t, err)

		v, err := s.SchemaVersion(ctx)
		assertNoErr(t, err)
		if v != 1 {
			t.Fatalf("expected version 1, got %d", v)
		}
	})
}

func TestStoreBook(t *testing.T) {
	dir := t.TempDir()
	s, err := storage.NewStore(path.Join(dir, "TestStoreBook.db"))
	assertNoErr(t, err)
	defer s.Close()
	ctx := context.Background()
	err = s.MigrateTo(ctx, 1)
	assertNoErr(t, err)

	t.Run("new book", func(t *testing.T) {
		b, err := pointsbook.NewBook("test-book-1")
		assertNoErr(t, err)

		err = s.StoreBook(ctx, b)
		assertNoErr(t, err)

		b, err = s.GetBook(ctx, "test-book-1")
		assertNoErr(t, err)
		assertBook(t, b, "test-book-1", 0, 0)
	})
	t.Run("book with trx", func(t *testing.T) {
		b, err := pointsbook.NewBook("test-book-2")
		assertNoErr(t, err)
		b.Add(5)
		b.Spend(2)

		err = s.StoreBook(ctx, b)
		assertNoErr(t, err)
		assertBook(t, b, "test-book-2", 3, 2)

		b, err = s.GetBook(ctx, "test-book-2")
		assertNoErr(t, err)
		assertBook(t, b, "test-book-2", 3, 2)
	})
	t.Run("multiple stores", func(t *testing.T) {
		b, err := pointsbook.NewBook("test-book-3")
		assertNoErr(t, err)
		b.Add(5)

		err = s.StoreBook(ctx, b)
		assertNoErr(t, err)
		assertBook(t, b, "test-book-3", 5, 1)

		b.Spend(2)

		err = s.StoreBook(ctx, b)
		assertNoErr(t, err)
		assertBook(t, b, "test-book-3", 3, 2)

		b, err = s.GetBook(ctx, "test-book-3")
		assertNoErr(t, err)
		assertBook(t, b, "test-book-3", 3, 2)
	})
}

func assertErr(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func assertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %s", err.Error())
	}
}

func assertBook(t *testing.T, b *pointsbook.Book, id string, points, trxLen int) {
	t.Helper()
	if b == nil {
		t.Fatal("expected non-nil book, got nil")
	}
	if b.ID() != id {
		t.Fatalf("expected book ID %s, got %s", strconv.Quote(id), strconv.Quote(b.ID()))
	}
	if b.CurrentPoints() != points {
		t.Fatalf("expected %d points, got %d", points, b.CurrentPoints())
	}
	if len(b.Transactions()) != trxLen {
		t.Fatalf("expected %d transactions, got %d", trxLen, len(b.Transactions()))
	}
	if len(b.UncommittedTransactions()) != 0 {
		t.Fatalf("expected 0 uncommitted transactions, got %d", len(b.UncommittedTransactions()))
	}
}
