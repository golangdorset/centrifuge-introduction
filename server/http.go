package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/centrifugal/centrifuge"
)

// authMiddleware sets Centrifuge credentials for the current request.
//
// Centrifuge expects all requests to have credentials attached. In this case
// we're not authorising users, so we just set a blank user ID.
func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cred := &centrifuge.Credentials{
			UserID: "",
		}
		newCtx := centrifuge.SetCredentials(ctx, cred)
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}

// broadcastHandler sends a message to all users subscribed to the broadcast
// channel.
func (a *app) broadcastHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if _, err = a.node.Publish("broadcast", b); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// publishHandler sends a message to users subscribed to their own user channel.
func (a *app) publishHandler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user == "" {
		http.Error(w, "missing user", http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if _, err = a.node.Publish(user, b); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// initHTTP initialises the HTTP server and router.
func (a *app) initHTTP() error {
	mux := http.NewServeMux()

	// Centrifuge endpoints.
	websocketHandler := centrifuge.NewWebsocketHandler(a.node, centrifuge.WebsocketConfig{})
	mux.Handle("/v1/connection/websocket", authMiddleware(websocketHandler))

	sockjsHandler := centrifuge.NewSockjsHandler(a.node, centrifuge.SockjsConfig{
		URL:                      "https://cdn.jsdelivr.net/npm/sockjs-client@1/dist/sockjs.min.js",
		HandlerPrefix:            "/connection/sockjs",
		WebsocketReadBufferSize:  1024,
		WebsocketWriteBufferSize: 1024,
	})
	mux.Handle("/v1/connection/sockjs", authMiddleware(sockjsHandler))

	// Message endpoints.
	mux.HandleFunc("/v1/message/broadcast", a.broadcastHandler)
	mux.HandleFunc("/v1/message/publish", a.publishHandler)

	a.h = &http.Server{
		Addr:    a.bind,
		Handler: mux,
	}

	return nil
}

// serveHTTP starts the HTTP server and sends an error to ch when finished.
func (a *app) serveHTTP(ch chan error) {
	ch <- a.h.ListenAndServe()
}
