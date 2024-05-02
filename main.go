package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

/*
dummy pool test
*/

type ping struct {
	init bool
	current int
	data chan struct {
		data int
	}
	finished chan int
}

func (g *ping) send(data int) {
	dataset := struct{ data int }{data}
	g.data <- dataset
}

func (g *ping) receive() {
	for {
		select {
			case v := <- g.data:
				g.current -= 1
				fmt.Println("receive:", v.data)
		}
	}
}

func main() {
	f, err := os.OpenFile("test.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}

	go func ()  {
		for i := range 100 {
			i := i
			go f.Write([]byte(fmt.Sprintf("%d\n",i)))
		}
	}()

	for i := range 200 {
		i := i
		go f.Write([]byte(fmt.Sprintf("%d\n",i)))
	}
	
	
	time.Sleep(time.Second * 20)
}