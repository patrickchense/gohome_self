package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var tk Config
var v = getToken()
var _ = json.Unmarshal([]byte(v), &tk)

func webhookGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Here is the token! ==> %s\n", tk.VerifyToken)

	token := tk.VerifyToken

	tokenTrue := r.URL.Query().Get("hub.verify_token")
	hubChallenge := r.URL.Query().Get("hub.challenge")

	if tokenTrue == token {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(hubChallenge))
	} else {
		fmt.Fprint(w, "Nay! Tokens don't match")
	}
}
