package core_test

import (
	"fmt"
	"testing"
)

func TestController2(t *testing.T) {
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
                c.current -= 1

                if c.started && c.current == 0 {
                    c.done <- 1
                    return
                }
            }
        }
    }()

    for i := range c.size {
        i := i
        c.current += 1
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