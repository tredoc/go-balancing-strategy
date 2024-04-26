package balancer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand/v2"
	"net/url"
	"sync"
	"testing"
)

func getNewRandomNode() *Node {
	port := rand.IntN(3000) + 1
	s, _ := url.Parse(fmt.Sprintf("http://localhost:%d", port))
	return NewNode(s, fmt.Sprintf("%d", port))
}

func TestBalancer_NextNode(t *testing.T) {
	t.Run("test strategy with random node select", func(t *testing.T) {
		s1 := getNewRandomNode()
		s2 := getNewRandomNode()
		s3 := getNewRandomNode()

		list := []*Node{s1, s2, s3}
		balancer := NewBalancer(list, NewRandomStrategy(int64(len(list))))

		nextNode := balancer.NextNode()
		assert.Equal(t, nextNode, s2, "second node")
	})

	t.Run("test strategy with single node to be selected", func(t *testing.T) {
		s1 := getNewRandomNode()

		list := []*Node{s1}
		balancer := NewBalancer(list, NewRandomStrategy(int64(len(list))))

		nextNode := balancer.NextNode()
		assert.Equal(t, nextNode, s1, "first node")
	})

	t.Run("test strategy with six node selected in pseudo-random way", func(t *testing.T) {
		s1 := getNewRandomNode()
		s2 := getNewRandomNode()
		s3 := getNewRandomNode()

		list := []*Node{s1, s2, s3}
		balancer := NewBalancer(list, NewRandomStrategy(int64(len(list))))

		nextNode := balancer.NextNode()
		assert.Equal(t, nextNode, s2, "second node")

		nextNode = balancer.NextNode()
		assert.Equal(t, nextNode, s3, "third node")

		nextNode = balancer.NextNode()
		assert.Equal(t, nextNode, s1, "first node")

		nextNode = balancer.NextNode()
		assert.Equal(t, nextNode, s1, "first node")

		nextNode = balancer.NextNode()
		assert.Equal(t, nextNode, s3, "third node")

		nextNode = balancer.NextNode()
		assert.Equal(t, nextNode, s1, "first node")
	})

	t.Run("test round robin strategy on six nodes", func(t *testing.T) {
		s1 := getNewRandomNode()
		s2 := getNewRandomNode()
		s3 := getNewRandomNode()

		list := []*Node{s1, s2, s3}
		balancer := NewBalancer(list, NewRoundRobinStrategy())
		counter := NewRoundRobinStrategy()

		var wg sync.WaitGroup
		count := 6
		wg.Add(count)

		for i := 0; i < count; i++ {
			go func() {
				defer wg.Done()
				nextNode := balancer.NextNode()
				c := counter.Next(len(list))
				assert.Equal(t, nextNode, list[c], i)
			}()
		}

		wg.Wait()
	})
}
