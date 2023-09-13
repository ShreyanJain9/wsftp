package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

func handleFTPConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Handle FTP commands and responses here
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		command := strings.TrimSpace(string(p))
		handleFTPCommand(conn, command)
	}
}

func handleFTPCommand(conn *websocket.Conn, command string) {
	// Implement command handling logic here
	// A switch statement to handle different FTP commands
	switch command {
	case "USER username":
		// Handle USER command for authentication
		// Respond with 331 User name okay, need password.
		// or 530 Not logged in as needed
		conn.WriteMessage(websocket.TextMessage, []byte("331 User name okay, need password."))
	case "PASS password":
		// Handle PASS command for password verification
		// Respond with 230 User logged in.
		// or 530 Not logged in as needed
		conn.WriteMessage(websocket.TextMessage, []byte("230 User logged in."))
	case "CWD directory":
		// Handle CWD command to change working directory
		// Respond as appropriate
		conn.WriteMessage(websocket.TextMessage, []byte("250 Directory changed successfully."))
	case "PWD":
		// Handle PWD command to print working directory
		// Respond with the current directory path
		conn.WriteMessage(websocket.TextMessage, []byte("257 \"/current/directory/path\""))
	case "LIST":
		// Handle LIST command for directory listing
		// Respond with the directory listing data
		conn.WriteMessage(websocket.TextMessage, []byte("150 Opening data connection for directory listing."))
		// Implement data channel logic to send directory listing
		// Respond with "226 Directory listing completed." when done
	case "RETR filename":
		// Handle RETR command for file retrieval
		// Respond with "150 Opening data connection for file retrieval."
		// Implement data channel logic to send the file
		// Respond with "226 Transfer complete." when done
	case "STOR filename":
		// Handle STOR command for file upload
		// Respond with "150 Opening data connection for file upload."
		// Implement data channel logic to receive and store the file
		// Respond with "226 Transfer complete." when done
	case "QUIT":
		// Handle QUIT command to disconnect
		conn.WriteMessage(websocket.TextMessage, []byte("221 Goodbye!"))
		conn.Close()
	default:
		// Handle unknown commands or send a "502 Command not implemented" response
		conn.WriteMessage(websocket.TextMessage, []byte("502 Command not implemented."))
	}
}

func main() {
	http.HandleFunc("/ftp", handleFTPConnection)
	fmt.Println("FTP server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
