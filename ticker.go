package main

import (
	"fmt"
	"net/http"
	"time"
)

type TickHandler struct {
	ticker          <-chan time.Time
	siskil          chan int
	actionBroadcast chan int
}

func (tHandler *TickHandler) startTicker(w http.ResponseWriter, req *http.Request) {
	tHandler.ticker = time.NewTicker(time.Second * 5).C
	tHandler.siskil = make(chan int, 1)
	fmt.Println("Starting recurring fetch operation")
	go tick(tHandler)
}

func (tHandler *TickHandler) tick(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Starting a single tick")
	tHandler.actionBroadcast <- 1
}

func (tHandler *TickHandler) stopTicker(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Stopping recurring operation")
	tHandler.siskil <- 1
}

func tick(tHandler *TickHandler) {
	for {
		select {
		case <-tHandler.ticker:
			tHandler.actionBroadcast <- 1
		case <-tHandler.siskil:
			return
		}
	}
}
