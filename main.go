package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/wonwooseo/sb-scanner/services"
)

var respBytes []byte
var lastUpdate time.Time

func main() {
	http.HandleFunc("/latest", ServeLatestCommitList)
	http.ListenAndServe(":80", nil)
}

// ServeLatestCommitList _
func ServeLatestCommitList(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET %s: %s", r.RequestURI, r.RemoteAddr)
	if respBytes == nil || time.Now().Sub(lastUpdate) >= time.Hour {
		var err error
		cacheCommitList, err := services.SearchLatestCommit()
		if err != nil {
			log.Fatal(err)
		}
		respBytes, err = json.Marshal(cacheCommitList)
		if err != nil {
			log.Println(err)
			return
		}
		lastUpdate = time.Now()
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(respBytes)
}
