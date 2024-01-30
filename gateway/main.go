package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	addr := flag.String("addr", ":6001", "the listen address of the gateway HTTP server")
	http.HandleFunc("/invoice", makeAPIFunc(handleGetInvoice))
	logrus.Infof("gateway HTTP server running on port %s", *addr)
	log.Fatal(http.ListenAndServe(":6001", nil))
}

func handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, map[string]any{
		"invoice": "some invoice",
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusOK, map[string]any{
				"error": err,
			})
		}
	}
}
