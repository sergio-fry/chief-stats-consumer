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
		client_id_cookie := clientCookie(r)
		w.Header().Add("Set-Cookie", client_id_cookie.String())

		session_id_cookie := sessionCookie(r)
		w.Header().Add("Set-Cookie", session_id_cookie.String())

		var job_params []interface{}

		job_params = append(job_params, "PageViewConsumerJob")
		job_params = append(job_params, []string{
			client_id_cookie.Value,
			session_id_cookie.Value,
			r.Header.Get("X-Real-IP"),
			r.UserAgent(),
			r.FormValue("url"),
			r.FormValue("referer"),
		})

		// Add a job to a queue
		// client_id, session_id, ip, user_agent, url, referer
		go workers.Enqueue("pageview", "JobRunner", job_params)

		log.Print(job_params)

		pic, err := os.Open("./public/blank_1px.gif")

		if err != nil {
			log.Fatal(err)
		}

		w.Header().Add("Cache-Control", "private, no-cache, no-store, must-revalidate, max-age=0")
		w.Header().Add("Pragma", "no-cache")
		w.Header().Add("Expires", "0")

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

func clientCookie(r *http.Request) *http.Cookie {
	cookie, _ := r.Cookie("cid")

	if cookie == nil {
		cookie = &http.Cookie{Name: "cid", Value: randomID()}
	}
	cookie.Expires = time.Now().Add(3 * Year)

	return cookie
}

func sessionCookie(r *http.Request) *http.Cookie {
	cookie, _ := r.Cookie("sid")

	if cookie == nil {
		cookie = &http.Cookie{Name: "sid", Value: randomID()}
	}

	cookie.Expires = time.Now().Add(30 * time.Minute)

	return cookie
}
