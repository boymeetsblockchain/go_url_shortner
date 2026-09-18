package main

import (
	"fmt"
	"net/http"

	"github.com/boymeetsblockchain/url_shortner/internal/api"
	"github.com/boymeetsblockchain/url_shortner/internal/store"
)

const addr = ":8080"

func main() {
	s := store.New()
	h := api.NewHandler(s, "http://localhost"+addr)

	fmt.Println("listening on", addr)
	http.ListenAndServe(addr, h.Routes())
}
