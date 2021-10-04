package main

import (
	"net/http"
)

func (app *Application) HandleMessages() {
	for {
		msg := <-app.Broadcast
		// app.InfoLog.Printf("msg: %+v\n", msg)
		for client := range app.Clients {
			if err := client.WriteJSON(msg); err != nil {
				msg := app.LeaveChat(client)
				go func() {
					app.Broadcast <- msg
				}()
			}
		}
	}
}

func (app *Application) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := app.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.ErrorLog.Fatal(err)
	}
	defer ws.Close()

	app.Clients[ws] = ""
	msg := app.CheckOnlineAndJoin()
	app.Broadcast <- msg

	for {
		msg := Message{}
		if err := ws.ReadJSON(&msg); err != nil {
			msg = app.LeaveChat(ws)
			app.Broadcast <- msg
			break
		}
		switch msg.Type {
		case TYPE_WS_JOIN:
			app.JoinChat(ws, &msg)
			app.Broadcast <- msg
		case TYPE_WS_MSG:
			app.Broadcast <- msg
		}
	}
}

func (app *Application) HandleCheckUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		m := map[string]interface{}{
			"status_code": http.StatusBadRequest,
			"username":    username,
			"is_allow":    false,
			"reason":      "Bad Request",
		}
		app.SendJSON(w, http.StatusOK, m)
		return
	}

	isExists := app.IsUsernameExists(username)
	if isExists {
		m := map[string]interface{}{
			"status_code": http.StatusOK,
			"username":    username,
			"is_allow":    false,
			"reason":      "Username already in use",
		}
		app.SendJSON(w, http.StatusOK, m)
		return
	}

	m := map[string]interface{}{
		"status_code": http.StatusOK,
		"username":    username,
		"is_allow":    true,
		"reason":      "",
	}
	app.SendJSON(w, http.StatusOK, m)
}

func (app *Application) HandleStaticFiles() http.Handler {
	return http.FileServer(app.FileSystem)
}
