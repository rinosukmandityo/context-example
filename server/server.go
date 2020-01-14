package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rinosukmandityo/context-example/google"
	"github.com/rinosukmandityo/context-example/userip"
)

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	// ctx is the Context for this handler. Calling cancel closes the
	// ctx.Done channel, which is the cancellation signal for requests
	// started by this handler.

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	timeout, e := time.ParseDuration(r.FormValue("timeout"))
	if e == nil {
		// The request has a timeout, so create a context that is
		// canceled automatically when the timeout expires.
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel() // Cancel ctx as soon as handleSearch returns.

	// Check the search query.
	q := r.FormValue("q")
	if q == "" {
		http.Error(w, "no query", http.StatusBadRequest)
		return
	}

	// Store the user IP in ctx for use by code in other packages.
	userIP, e := userip.FromRequest(r)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
	}
	ctx = userip.NewContext(ctx, userIP)

	// Run the Google search and print the results.
	start := time.Now()
	results, e := google.Search(ctx, q)
	elapsed := time.Since(start)

	fmt.Println(results, elapsed)

}
