package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func (app *Application) AmountOnlineAndJoin() (int, int) {
	join := 0
	for _, username := range app.Clients {
		if username != "" {
			join++
		}
	}
	return len(app.Clients), join
}

func (app *Application) IsUsernameExists(username string) bool {
	for _, uname := range app.Clients {
		if uname == username {
			return true
		}
	}
	return false
}

func (app *Application) JoinChat(ws *websocket.Conn, msg *Message) {
	app.Clients[ws] = msg.Username
	cOn, cJoin := app.AmountOnlineAndJoin()
	msg.AmountOnline = cOn
	msg.AmountJoin = cJoin
	app.InfoLog.Output(2, fmt.Sprintf("%s join chat\n", app.Clients[ws]))
}

func (app *Application) LeaveChat(ws *websocket.Conn) Message {
	msg := Message{}
	username := app.Clients[ws]
	delete(app.Clients, ws)

	defer func() {
		app.InfoLog.Output(3, fmt.Sprintf("%s leave chat\n", username))
		ws.Close()
	}()

	cOn, cJoin := app.AmountOnlineAndJoin()
	if username != "" {
		msg = Message{
			Type:         TYPE_WS_LEAVE,
			Username:     username,
			AmountOnline: cOn,
			AmountJoin:   cJoin,
		}
	} else {
		username = ws.RemoteAddr().String()
		msg = Message{
			Type:         TYPE_WS_CHECK_ONLINE_JOIN,
			AmountOnline: cOn,
			AmountJoin:   cJoin,
		}
	}
	return msg
}

func (app *Application) CheckOnlineAndJoin() Message {
	cOn, cJoin := app.AmountOnlineAndJoin()
	msg := Message{
		Type:         TYPE_WS_CHECK_ONLINE_JOIN,
		AmountOnline: cOn,
		AmountJoin:   cJoin,
	}
	return msg
}

func (app *Application) SendJSON(w http.ResponseWriter, statusCode int, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	vJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Write(vJson)
}
