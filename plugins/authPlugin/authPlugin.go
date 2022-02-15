package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const pluginName = "authPlugin"

func main() {}

func init() {
	fmt.Printf("%s plugin is loaded!\n", pluginName)
}

// HandlerRegisterer is the name of the symbol krakend looks up to try and register plugins
var HandlerRegisterer registerer = registerer(pluginName)

type registerer string

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(
		context.Context,
		map[string]interface{},
		http.Handler) (http.Handler, error),
)) {
	fmt.Printf("%s registered!\n", pluginName)
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(ctx context.Context, extra map[string]interface{}, handler http.Handler) (http.Handler, error) {

	client := &http.Client{Timeout: 3 * time.Second}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/v18/bootstrap" {
			data := url.Values{}
			data.Set("product_line_name", "remitly")

			rq, err := http.NewRequest(http.MethodPost, "http://auth-service-preprod.us-west-2.remitly.com/v1/session/validate", nil)
			rq.Header.Add("Authorization", "Bearer "+req.Header["Authorization"][0])

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			rq.Header.Set("Content-Type", "application/json")

			rs, err := client.Do(rq)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotAcceptable)
				return
			}
			defer rs.Body.Close()

			if rs.StatusCode >= 300 {
				w.WriteHeader(rs.StatusCode)
			} else {
				r2 := new(http.Request)
				*r2 = *req

				handler.ServeHTTP(w, r2)
			}
		} else {
			r2 := new(http.Request)
			*r2 = *req

			handler.ServeHTTP(w, r2)
		}

	}), nil
}
