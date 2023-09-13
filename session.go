package main

import (
	"fmt"
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
var cwd = map[*websocket.Conn]string{}

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

func isAuthenticated(conn *websocket.Conn) bool {
	ftpContext := getFTPContext(conn)
	if ftpContext != nil {
		return ftpContext.IsAuthenticated
	}
	return false
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

func handleCWD(conn *websocket.Conn, args string) {
	ftpContext := getFTPContext(conn)
	if ftpContext == nil || !ftpContext.IsAuthenticated {
		sendResponse(conn, 530, NotLoggedIn)
		return
	}

	newDir := strings.TrimSpace(args)
	// Implement logic to validate and change the directory
	// Set the current working directory in the context
	ftpContext.WorkingDir = newDir
	sendResponse(conn, 250, "Directory changed successfully.")
}

func handlePWD(conn *websocket.Conn) {
	ftpContext := getFTPContext(conn)
	if ftpContext == nil || !ftpContext.IsAuthenticated {
		sendResponse(conn, 530, NotLoggedIn)
		return
	}

	currentDir := ftpContext.WorkingDir
	sendResponse(conn, 257, "\""+currentDir+"\"")
}

func handleLIST(conn *websocket.Conn) {
	ftpContext := getFTPContext(conn)
	if ftpContext == nil || !ftpContext.IsAuthenticated {
		sendResponse(conn, 530, NotLoggedIn)
		return
	}

	// Implement logic to generate directory listing data
	sendResponse(conn, 150, "Opening data connection for directory listing.")
	// Implement data channel logic to send directory listing
	// Respond with "226 Directory listing completed." when done
	sendResponse(conn, 226, "Directory listing completed.")
}

func handleRETR(conn *websocket.Conn, args string) {
	ftpContext := getFTPContext(conn)
	if ftpContext == nil || !ftpContext.IsAuthenticated {
		sendResponse(conn, 530, NotLoggedIn)
		return
	}

	filename := strings.TrimSpace(args)
	// Implement logic to retrieve and send the specified file
	sendResponse(conn, 150, "Opening data connection for file retrieval.")
	// Implement data channel logic to send the file
	// Respond with "226 Transfer complete." when done

	fmt.Println("filename: ", filename) // to make compiler happy

	sendResponse(conn, 226, "Transfer complete.")
}

func handleSTOR(conn *websocket.Conn, args string) {
	ftpContext := getFTPContext(conn)
	if ftpContext == nil || !ftpContext.IsAuthenticated {
		sendResponse(conn, 530, NotLoggedIn)
		return
	}

	filename := strings.TrimSpace(args)
	// Implement logic to receive and store the uploaded file
	sendResponse(conn, 150, "Opening data connection for file upload.")
	// Implement data channel logic to receive and store the file
	// Respond with "226 Transfer complete." when done

	fmt.Println("filename: ", filename) // to make compiler happy

	sendResponse(conn, 226, "Transfer complete.")
}

func handleQUIT(conn *websocket.Conn) {
	connContextMap.Delete(conn) // Remove the connection from the map
	sendResponse(conn, 221, "Goodbye!")
	conn.Close()
}
