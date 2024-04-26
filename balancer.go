package balancer

import (
	"math/rand"
	"net/url"
	"sync/atomic"
)

type Node struct {
	Name    string
	Address *url.URL
}

func NewNode(address *url.URL, name string) *Node {
	return &Node{Address: address, Name: name}
}

type BalancingStrategy interface {
	Next(int) int
}

type Balancer struct {
	nodes    []*Node
	strategy BalancingStrategy
}

func NewBalancer(nodes []*Node, strategy BalancingStrategy) *Balancer {
	return &Balancer{nodes: nodes, strategy: strategy}
}

func (b Balancer) NextNode() *Node {
	idx := b.strategy.Next(len(b.nodes))
	return b.nodes[idx]
}

type RandomStrategy struct {
	rand *rand.Rand
}

func NewRandomStrategy(seed int64) *RandomStrategy {
	s := rand.NewSource(seed)
	r := rand.New(s)
	return &RandomStrategy{rand: r}
}

func (r RandomStrategy) Next(length int) int {
	return r.rand.Intn(length)
}

type RoundRobinStrategy struct {
	currentIdx atomic.Int64
}

func NewRoundRobinStrategy() *RoundRobinStrategy {
	return &RoundRobinStrategy{}
}

func (rr *RoundRobinStrategy) Next(length int) int {
	current := rr.currentIdx.Load()
	if int(current) >= length-1 {
		rr.currentIdx.CompareAndSwap(current, 0)
		return int(rr.currentIdx.Load())
	}

	rr.currentIdx.CompareAndSwap(current, current+1)
	return int(rr.currentIdx.Load())
}
