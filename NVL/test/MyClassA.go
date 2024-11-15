package main

import (
	"sync"
	"sync/atomic"
)

type MyClassA struct {
	Name        string
	Age         int
	IsSomething bool
	X, Y, Z     float64

	wg        sync.WaitGroup
	ba        []byte
	neighbour *MyClassA
	counter   atomic.Int32
	m         *sync.Mutex
}

func (mca *MyClassA) CanDoSomething() bool {
	return mca.IsSomething
}

func FromAny(x any) (mca *MyClassA) {
	return x.(*MyClassA)
}
