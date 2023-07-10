package kademlia

import (
	"math/big"
)

type bucket struct {
	contacts bigIntMap[Contact]
	from, to *big.Int
}

// cmp returns an integer denoting z's inclusion into the keyspace. The result
// will be 0 if z is contained in the keyspace, -1 if z is lower than the lowest
// possible key in the keyspace, and +1 if z is larger than the biggest
// possible key in the keyspace.
func (b *bucket) cmp(z *big.Int) int {
	x, y := z.Cmp(b.from), z.Cmp(b.to)
	switch {
	case x < 0:
		return -1
	case x >= 0 && y <= 0:
		return 0
	default:
		return 1
	}
}

func (b *bucket) put(c Contact, k int) (ok bool) {
	if b.contacts.len() == 0 {
		b.contacts.set(c.ID(), c)
		return true
	}

	if b.contacts.contains(c.ID()) {
		return true
	}

	if k < 0 || b.contacts.len() < k {
		b.contacts.set(c.ID(), c)
		return true
	}

	return false
}

func (b *bucket) forEach(f func(Contact) bool) {
	b.contacts.forEach(func(_ *big.Int, contact Contact) bool {
		return f(contact)
	})
}

var (
	zero = big.NewInt(0)
	one  = big.NewInt(1)
	two  = big.NewInt(2)
)

func (b *bucket) split() (lo, hi *bucket, ok bool) {
	var x, y big.Int

	x.Add(b.from, b.to)
	x.Div(&x, two)
	if x.Cmp(zero) == 0 {
		return
	}

	y.Add(&x, one)
	if y.Cmp(b.to) >= 0 {
		return
	}

	lo = newBucket(b.from, &x)
	hi = newBucket(&y, b.to)
	ok = true

	b.forEach(func(c Contact) bool {
		switch {
		case hi.cmp(c.ID()) == 0:
			hi.contacts.set(c.ID(), c)
		default:
			lo.contacts.set(c.ID(), c)
		}
		return true
	})
	return
}

func (b *bucket) depth() int {
	const (
		bits = 8
	)

	switch b.contacts.len() {
	case 0:
		return 0
	case 1:
		return bitAlign(b.contacts.first().ID().BitLen(), bits)
	}

	var min, max *big.Int
	var i int
	b.contacts.forEach(func(_ *big.Int, contact Contact) bool {
		id := contact.ID()

		switch i {
		case 0:
			min, max = id, id
		default:
			switch {
			case id.Cmp(min) < 0:
				min = id
			case id.Cmp(max) > 0:
				max = id
			}
		}
		i++

		return true
	})
	return bitAlign(max.BitLen(), bits) - new(big.Int).Xor(min, max).BitLen()
}

func newBucket(from, to *big.Int) *bucket {
	return &bucket{from: from, to: to}
}
