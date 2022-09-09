package kademlia

import (
	"math/big"
	"sort"
	"sync"
)

const (
	DefaultSize       = 160
	DefaultBucketSize = 8
	DefaultSplitLevel = 5
)

type Table interface {
	Put(contact Contact)
	ClosestContacts(id *big.Int, n int) []Contact
}

// table is a dynamic sized Kademlia routing table.
type table struct {
	mu      sync.RWMutex
	homeIDs intMap[struct{}]
	buckets []*bucket
	B       *big.Int
	k, b    int
}

func (t *table) containsHomeID(b *bucket) (ok bool) {
	ok = t.homeIDs.length() > 0
	t.homeIDs.forEach(func(id *big.Int, _ struct{}) bool {
		if b.cmp(id) != 0 {
			ok = false
			return false
		}
		return true
	})
	return
}

func (t *table) search(id *big.Int) (ret int) {
	ret = sort.Search(len(t.buckets), func(i int) bool {
		return t.buckets[i].cmp(id) < 1
	})
	if ret == len(t.buckets) {
		ret = -1
	}
	return
}

func (t *table) Put(contact Contact) {
	// Ignore contacts with ID larger than our ID space or with
	// the same ID as ours.
	if contact.ID().Cmp(t.B) > 0 || t.homeIDs.contains(contact.ID()) {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.buckets == nil {
		t.buckets = []*bucket{newBucket(zero, t.B)}
	}

	i := t.search(contact.ID())
	if i < 0 {
		return
	}

	b := t.buckets[i]
	if !b.put(contact, t.k) {
		// Try to split bucket if the bucket's range contains the node's
		// own ID or the depth d of the k-bucket in the routing tree
		// satisfies d ≢ 0 (mod b).
		if t.containsHomeID(b) || b.depth()%t.b != 0 {
			if lo, hi, ok := b.split(); ok {
				b = lo
				if b.cmp(contact.ID()) > 0 {
					b = hi
				}

				// If contact can be succesfully put into one of the recently
				// split buckets, commit changes.
				if b.put(contact, t.k) {
					t.buckets = append(t.buckets, nil)
					copy(t.buckets[i+2:], t.buckets[i+1:])
					t.buckets[i] = lo
					t.buckets[i+1] = hi
					return
				}
			}
		}
	}
}

func (t *table) closestContacts(id *big.Int, n int, f func(Contact)) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	i := t.search(id)
	if i < 0 {
		return
	}

	g := func(x Contact) bool {
		f(x)
		n--
		return true
	}

	// Buckets are traversed in ascending and descending order alternately.
	for j := i + 1; i >= 0 || j < len(t.buckets); i, j = i-1, j+1 {
		if i >= 0 {
			t.buckets[i].forEach(g)
		}
		if j < len(t.buckets) {
			t.buckets[j].forEach(g)
		}

		// Have we visited enough contacts?
		if n <= 0 {
			return
		}
	}
}

func (t *table) ClosestContacts(id *big.Int, n int) (ret []Contact) {
	ret = make([]Contact, 0, n)
	t.closestContacts(id, n, func(x Contact) {
		ret = append(ret, x)
	})

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].ID().Cmp(ret[j].ID()) < 1
	})
	if len(ret) > n {
		ret = ret[:n]
	}
	return
}

func (t *table) PutHomeID(id *big.Int) {
	t.mu.Lock()
	t.homeIDs.set(id, struct{}{})
	t.mu.Unlock()
}

func WithHomeID(id *big.Int) func(*table) {
	return func(t *table) {
		t.homeIDs.set(id, struct{}{})
	}
}

func WithSize(B int) func(*table) {
	return func(t *table) {
		t.B = bitMask(B)
	}
}

func WithBucketSize(k int) func(*table) {
	return func(t *table) {
		t.k = k
	}
}

func WithSplitLevel(b int) func(*table) {
	return func(t *table) {
		t.b = b
	}
}

func NewTable(opts ...func(*table)) Table {
	ret := table{
		B: bitMask(DefaultSize),
		k: DefaultBucketSize,
		b: DefaultSplitLevel,
	}
	for _, f := range opts {
		f(&ret)
	}
	return &ret
}
