package main

import (
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

const NotLoggedIn = "User not logged in."

// Defines a custom context structure for managing user, directory information, and authentication state.
type FTPContext struct {
	User            string
	WorkingDir      string
	IsAuthenticated bool
}

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

var cwdMutex sync.Mutex

var connContextMap sync.Map // Map to store context information for each connection

func setFTPContext(conn *websocket.Conn, ftpContext *FTPContext) {
	connContextMap.Store(conn, ftpContext)
}

func getFTPContext(conn *websocket.Conn) *FTPContext {
	val, ok := connContextMap.Load(conn)
	if !ok {
		return nil
	}
	context, _ := val.(*FTPContext)
	return context
}

func authenticate(conn *websocket.Conn, username, password string) bool {
	storedPassword, exists := users[username]
	if exists && storedPassword == password {
		ftpContext := getFTPContext(conn)
		ftpContext.IsAuthenticated = true
		return true
	}
	return false
}

func sendResponse(conn *websocket.Conn, code int, response string) {
	conn.WriteMessage(websocket.TextMessage, []byte(strings.Join([]string{strconv.Itoa(code), response}, " ")))
}

func sendData(conn *websocket.Conn, data string) {
	conn.WriteMessage(websocket.TextMessage, []byte(data))
}

func handleUSER(conn *websocket.Conn, args string) {
	username := strings.TrimSpace(args)
	if _, exists := users[username]; exists {
		ftpContext := &FTPContext{
			User: username,
		}
		setFTPContext(conn, ftpContext)
		sendResponse(conn, 331, "User name okay, need password.")
	} else {
		sendResponse(conn, 530, NotLoggedIn)
	}
}

func handlePASS(conn *websocket.Conn, args string) {
	password := strings.TrimSpace(args)
	ftpContext := getFTPContext(conn)
	username := ftpContext.User
	if authenticate(conn, username, password) {
		sendResponse(conn, 230, "User logged in.")
	} else {
		sendResponse(conn, 530, NotLoggedIn)
	}
}

func handleQUIT(conn *websocket.Conn) {
	connContextMap.Delete(conn) // Remove the connection from the map
	sendResponse(conn, 221, "Goodbye!")
	conn.Close()
}
