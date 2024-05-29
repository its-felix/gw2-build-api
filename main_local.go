//go:build !lambda

package main

import (
	"net/http"
)

type buildApiResponse struct {
	Id int `json:"id"`
}

func main() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: setupMux(),
	}
	_ = srv.ListenAndServe()
}
