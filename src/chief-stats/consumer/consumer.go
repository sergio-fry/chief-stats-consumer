package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/jrallison/go-workers"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const Year = 365 * 24 * time.Hour

func main() {
	workers.Configure(map[string]string{
		// location of redis instance
		"server": "localhost:6379",
		// instance of the database
		"database": "0",
		// number of connections to keep open with redis
		"pool": "30",
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process":   "1",
		"namespace": "chiefstats",
	})

	http.HandleFunc("/api/p/v", func(w http.ResponseWriter, r *http.Request) {
		//===========================================================================
		// Client id
		client_id_cookie, _ := r.Cookie("cid")

		if client_id_cookie == nil {
			client_id_cookie = &http.Cookie{Name: "cid", Value: randomID()}
		}

		client_id_cookie.Expires = time.Now().Add(3 * Year)

		w.Header().Add("Set-Cookie", client_id_cookie.String())

		//===========================================================================
		// Client id

		session_id_cookie, _ := r.Cookie("sid")

		if session_id_cookie == nil {
			session_id_cookie = &http.Cookie{Name: "sid", Value: randomID()}
		}

		session_id_cookie.Expires = time.Now().Add(30 * time.Minute)

		w.Header().Add("Set-Cookie", session_id_cookie.String())

		//===========================================================================

		// Add a job to a queue
		// client_id, session_id, ip, user_agent, url, referer
		go workers.Enqueue("pageview", "Add", []string{client_id_cookie.Value, session_id_cookie.Value})
		log.Printf("cid: %s sid: %s", client_id_cookie.Value, session_id_cookie.Value)

		pic, err := os.Open("./public/blank_1px.gif")

		if err != nil {
			log.Fatal(err)
		}

		io.Copy(w, pic)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func randomID() (id string) {
	data := make([]byte, 10)

	if _, err := rand.Read(data); err == nil {
		id = fmt.Sprintf("%x", sha256.Sum256(data))
	}

	return id
}
