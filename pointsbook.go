package pointsbook

import "errors"

type Book struct {
	id   string
	trxs []int
}

func NewBook(id string) (*Book, error) {
	if id == "" {
		return nil, errors.New("empty ID")
	}
	return &Book{id: id}, nil
}

func (b *Book) ID() string { return b.id }

func (b *Book) CurrentPoints() int {
	var p int
	for _, t := range b.trxs {
		p += t
	}
	return p
}

func (b *Book) Add(p uint) {
	b.trxs = append(b.trxs, int(p))
}

func (b *Book) Transactions() []int { return b.trxs }

func (b *Book) Spend(p uint) (ok bool) {
	pi := int(p)
	if b.CurrentPoints() < pi {
		return false
	}
	b.trxs = append(b.trxs, -1*pi)
	return true
}
