package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const outputHeaderName = "X-Friend-User"

func main() {}

func init() {
	fmt.Println("headerModPlugin plugin is loaded!")
}

// HandlerRegisterer is the name of the symbol krakend looks up to try and register plugins
var HandlerRegisterer registerer = registerer("headerModPlugin")

type registerer string

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(
		context.Context,
		map[string]interface{},
		http.Handler) (http.Handler, error),
)) {
	fmt.Println("headerModPlugin registered!")
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(ctx context.Context, extra map[string]interface{}, handler http.Handler) (http.Handler, error) {
	attachUserID, ok := extra["attachuserid"].(string)
	if !ok {
		panic(errors.New("incorrect config").Error())
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r2 := new(http.Request)
		*r2 = *r

		r2.Header.Set(outputHeaderName, string(attachUserID))

		handler.ServeHTTP(w, r2)
	}), nil
}

// func (r registerer) registerHandlers(ctx context.Context, extra map[string]interface{}, _ http.Handler) (http.Handler, error) {
// 	// check the passed configuration and initialize the plugin
// 	name, ok := extra["name"].(string)
// 	if !ok {
// 		return nil, errors.New("wrong config")
// 	}
// 	if name != string(r) {
// 		return nil, fmt.Errorf("unknown register %s", name)
// 	}
// 	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http handler
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))
// 	}), nil
// }
