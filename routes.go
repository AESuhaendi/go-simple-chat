package main

import "net/http"

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	// Websockets
	mux.HandleFunc("/ws", app.HandleConnections)

	// Check Username
	mux.HandleFunc("/check-username", app.HandleCheckUsername)

	// Static Files
	mux.Handle("/", app.HandleStaticFiles())

	return mux
}
