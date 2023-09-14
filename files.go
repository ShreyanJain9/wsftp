package main

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

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
