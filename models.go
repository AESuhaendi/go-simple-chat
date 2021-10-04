package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type         string `json:"type"`
	Username     string `json:"username"`
	Message      string `json:"message"`
	AmountOnline int    `json:"amount_online"`
	AmountJoin   int    `json:"amount_join"`
}

type Application struct {
	Upgrader   websocket.Upgrader
	Clients    map[*websocket.Conn]string
	Broadcast  chan Message
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
	FileSystem http.FileSystem
}
