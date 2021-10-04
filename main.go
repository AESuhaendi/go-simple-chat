package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	TYPE_WS_JOIN              = "TYPE_WS_JOIN"
	TYPE_WS_LEAVE             = "TYPE_WS_LEAVE"
	TYPE_WS_MSG               = "TYPE_WS_MSG"
	TYPE_WS_CHECK_ONLINE_JOIN = "TYPE_WS_CHECK_ONLINE_JOIN"
)

//go:embed public
var staticfiles embed.FS

func main() {
	var port string
	flag.StringVar(&port, "port", ":8000", "HTTP Server Port")
	flag.Parse()

	infoLog := log.New(os.Stdout, "[LOG] ", log.Lshortfile|log.LstdFlags)
	errorLog := log.New(os.Stderr, "[ERR] ", log.Lshortfile|log.LstdFlags)

	// Makes sure you put folder name correctly and exists
	fsStatic, err := fs.Sub(staticfiles, "public")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &Application{
		Upgrader:   websocket.Upgrader{},
		Clients:    make(map[*websocket.Conn]string),
		Broadcast:  make(chan Message),
		InfoLog:    infoLog,
		ErrorLog:   errorLog,
		FileSystem: http.FS(fsStatic),
	}

	go app.HandleMessages()

	srv := http.Server{
		Addr:         port,
		IdleTimeout:  2 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      app.routes(),
	}

	infoLog.Printf("Listening on: http://localhost%s", port)
	if err := srv.ListenAndServe(); err != nil {
		errorLog.Fatal(err)
	}
}
