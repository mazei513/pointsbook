package storage

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/mazei513/pointsbook"
	sqlite "modernc.org/sqlite"
)

var (
	ErrUninitialized = errors.New("store uninitialized")
)

type Store struct {
	db *sql.DB
}

func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`pragma foreign_keys = on;`)
	if err != nil {
		return nil, err
	}
	return &Store{db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) SchemaVersion(ctx context.Context) (int, error) {
	r := s.db.QueryRowContext(ctx, "select v from dbver;")
	var v int
	err := r.Scan(&v)
	if err != nil {
		sqliteErr := &sqlite.Error{}
		if errors.As(err, &sqliteErr) && strings.Contains(sqliteErr.Error(), "no such table: dbver") {
			return 0, ErrUninitialized
		}
		return 0, err
	}
	return v, nil
}

func (s *Store) MigrateTo(ctx context.Context, target int) error {
	cur, err := s.SchemaVersion(ctx)
	if errors.Is(err, ErrUninitialized) {
		cur = -1
	}
	if cur == target {
		return nil
	}

	if target < cur {
		return errors.New("migration target is lower than current version")
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	for i := cur; i < target; i++ {
		s, ok := migrateScripts[i+1]
		if !ok {
			tx.Rollback()
			return errors.New("target migration version does not exist")
		}
		_, err := tx.ExecContext(ctx, s)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.ExecContext(ctx, `update dbver set v=?;`, i+1)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *Store) StoreBook(ctx context.Context, b *pointsbook.Book) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `insert into books values (?) on conflict do nothing;`, b.ID())
	if err != nil {
		tx.Rollback()
		return err
	}

	btrxs := b.Transactions()
	start := len(btrxs) - len(b.UncommittedTransactions())
	for i := start; i < len(btrxs); i++ {
		_, err := tx.ExecContext(ctx, `insert into book_trxs values (?, ?, ?) ;`, b.ID(), i, btrxs[i])
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	b.CommitTransactions()

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *Store) GetBook(ctx context.Context, id string) (*pointsbook.Book, error) {
	r := s.db.QueryRowContext(ctx, `select count(*) from books where uid=?;`, id)
	var i int
	err := r.Scan(&i)
	if err != nil {
		return nil, err
	}
	if i != 1 {
		return nil, err
	}

	var btrxs []int
	rs, err := s.db.QueryContext(ctx, `select amount from book_trxs where uid=? order by trx_idx asc;`, id)
	if err != nil {
		return nil, err
	}
	for rs.Next() {
		var i int
		err := rs.Scan(&i)
		if err != nil {
			return nil, err
		}
		btrxs = append(btrxs, i)
	}

	return pointsbook.BookFromTransactions(id, btrxs)
}

var migrateScripts = map[int]string{
	0: `create table dbver(v integer not null); insert into dbver values (-1);`,
	1: `
create table books(uid text primary key not null);
create table book_trxs(
	uid text not null,
	trx_idx int not null,
	amount int not null,
	primary key(uid, trx_idx),
	foreign key(uid) references books(uid) on delete restrict
);
`,
}
