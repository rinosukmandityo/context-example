package google

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rinosukmandityo/context-example/userip"
)

type Results []Result

type Result struct {
	Title, URL string
}

func Search(ctx context.Context, query string) (Results, error) {
	// Prepare the Google Search API request.
	var results Results
	r, e := http.NewRequest("GET", "https://ajax.googleapis.com/ajax/services/search/web?v=1.0", nil)
	if e != nil {
		return results, e
	}
	q := r.URL.Query()
	q.Set("q", query)

	// If ctx is carrying the user IP address, forward it to the server.
	// Google APIs use the user IP to distinguish server-initiated requests
	// from end-user requests.
	if netIP, ok := userip.FromContext(ctx); ok {
		q.Set("userip", netIP.String())
	}
	r.URL.RawQuery = q.Encode()

	e = httpDo(ctx, r, func(resp *http.Response, e error) error {
		if e != nil {
			return e
		}

		defer resp.Body.Close()

		// Parse the JSON search result.
		// https://developers.google.com/web-search/docs/#fonje

		var data struct {
			ResponseData struct {
				Results []struct {
					TitleNoFormatting string
					URL               string
				}
			}
		}

		if e := json.NewDecoder(resp.Body).Decode(&data); e != nil {
			return e
		}

		for _, res := range data.ResponseData.Results {
			results = append(results, Result{res.TitleNoFormatting, res.URL})
		}

		return nil
	})
	// httpDo waits for the closure we provided to return, so it's safe to
	// read results here.

	return results, nil
}

func httpDo(ctx context.Context, r *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.

	c := make(chan error, 1)
	r = r.WithContext(ctx)
	go func() { c <- f(http.DefaultClient.Do(r)) }()

	select {
	case <-ctx.Done():
		<-c // Wait for f to return
		return ctx.Err()
	case e := <-c:
		return e
	}
}
