/*
just dummy testing
*/
package core_test

import (
	"fmt"
	"testing"
)

type Controller struct {
	started  bool
	size     uint32
	current  uint32
	pingpong chan struct {
		message string
	}
	done chan int
}

func TestController1(t *testing.T) {
	c := &Controller{
		size:     400,
		pingpong: make(chan struct{ message string }, 50),
		done:     make(chan int),
	}

	go func() {
		defer close(c.pingpong)
		for {
			select {
			case v := <-c.pingpong:
				t.Log(v)
				c.current += 1

				if c.started && c.current == c.size {
					c.done <- 1
					return
				}
			}
		}
	}()

	for i := range c.size {
		i := i
		go func() {
			if i == 0 {
				c.started = true
			}
			c.pingpong <- struct{ message string }{message: fmt.Sprintf("hello thread %d", i)}
		}()
	}

	<-c.done
	t.Log("done")
}
