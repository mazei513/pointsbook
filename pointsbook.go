package pointsbook

type Book struct {
	id   string
	trxs []int
}

func NewBook(id string) *Book {
	return &Book{id: id}
}

func (b *Book) ID() string { return b.id }

func (b *Book) CurrentPoints() int {
	var p int
	for _, t := range b.trxs {
		p += t
	}
	return p
}

func (b *Book) AddPoints(p int) {
	b.trxs = append(b.trxs, p)
}

func (b *Book) Transactions() []int { return b.trxs }
