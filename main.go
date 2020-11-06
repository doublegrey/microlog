package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/doublegrey/microlog/utils"
)

func validateToken(handler func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Token")
		if len(strings.TrimSpace(token)) < 1 {
			w.WriteHeader(401)
			fmt.Fprintf(w, "Access Token does not exists or is invalid")
		} else {
			handler(w, r)
		}
	})
}

func main() {
	err := utils.Config.Parse()
	if err != nil {
		log.Fatalf("Failed to parse config: %v\n", err)
	}

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", utils.Config.Host, utils.Config.Port), nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
