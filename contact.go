package kademlia

import "math/big"

type Contact interface {
	ID() *big.Int
}
