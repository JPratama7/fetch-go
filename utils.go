package main

import (
	"net/http"

	"github.com/enetx/surf"
)

func NewClient() *http.Client {
	return surf.NewClient().
		Builder().
		Singleton().
		AddHeaders("Content-Encoding", "gzip").
		DNSOverTLS().Cloudflare().
		Impersonate().
		RandomOS().
		Chrome().
		Build().
		Std()
}
