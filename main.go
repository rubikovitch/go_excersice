package main

import (
	"log"
	"net/http"
)

func main() {
	tickHandler := TickHandler{}
	tickHandler.actionBroadcast = make(chan int)
	listenAndHandle(tickHandler.actionBroadcast)
	http.HandleFunc("/start", tickHandler.startTicker)
	http.HandleFunc("/tick", tickHandler.tick)
	http.HandleFunc("/kill", tickHandler.stopTicker)
	http.HandleFunc("/fetchRecent", localRecentTimeHandler)
	http.HandleFunc("/fetchAll", localAllTimeHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
