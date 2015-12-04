package main

import (
	"github.com/jrallison/go-workers"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	workers.Configure(map[string]string{
		// location of redis instance
		"server": "localhost:6379",
		// instance of the database
		"database": "0",
		// number of connections to keep open with redis
		"pool": "30",
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process": "1",
	})

	http.HandleFunc("/api/p/v", func(w http.ResponseWriter, r *http.Request) {

		pic, err := os.Open("./public/blank_1px.gif")

		if err != nil {
			log.Fatal(err)
		}

		io.Copy(w, pic)
	})

	go workers.StatsServer(8081)
	go workers.Run()
	log.Fatal(http.ListenAndServe(":8080", nil))

}
