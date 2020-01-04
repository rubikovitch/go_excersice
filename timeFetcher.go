package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type SimpleTime struct {
	Datetime string
}

func listenAndHandle(sigStart chan int) {
	go func() {
		for {
			select {
			case _, ok := <-sigStart:
				if !ok {
					break
				}
				time := fetchTime()
				saveTimeToDB(time)
			}
		}
	}()
}

func fetchTime() string {
	log.Println("Fetching time now")
	resp, err := http.Get("http://worldtimeapi.org/api/timezone/Europe/Minsk")
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	defer resp.Body.Close()
	var simpleTime SimpleTime
	jsonErr := json.NewDecoder(resp.Body).Decode(&simpleTime)
	if jsonErr != nil {
		log.Fatalln(err)
		return ""
	}
	log.Printf("Got time: %+v \n", simpleTime.Datetime)
	return simpleTime.Datetime
}
