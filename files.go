package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

func handleRETR(conn *websocket.Conn, args string) {
	ftpContext := getFTPContext(conn)
	if ftpContext == nil || !ftpContext.IsAuthenticated {
		sendResponse(conn, 530, NotLoggedIn)
		return
	}

	filename := filepath.Join(Dir(ftpContext), strings.TrimSpace(args))
	// Implement logic to retrieve and send the specified file
	sendResponse(conn, 150, "Opening data connection for file retrieval.")
	// Implement data channel logic to send the file

	file, err := os.Open(filename)
	if err != nil {
		sendResponse(conn, 550, "File not found.")
		return
	}
	defer file.Close()

	// Read and send the file data over the WebSocket connection
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			sendResponse(conn, 550, "Failed to read file.")
			return
		}
		conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
	}

	// Respond with "226 Transfer complete." when done

	sendResponse(conn, 226, "Transfer complete.")
}

func handleSTOR(conn *websocket.Conn, args string) {
	ftpContext := getFTPContext(conn)
	if ftpContext == nil || !ftpContext.IsAuthenticated {
		sendResponse(conn, 530, NotLoggedIn)
		return
	}

	filename := filepath.Join(Dir(ftpContext), strings.TrimSpace(args))
	// Implement logic to receive and store the uploaded file
	sendResponse(conn, 150, "Opening data connection for file upload.")
	// Implement data channel logic to receive and store the file

	file, err := os.Create(filename)
	if err != nil {
		sendResponse(conn, 550, "Failed to create file.")
		return
	}
	defer file.Close()

	// Receive and store the file data from the WebSocket connection
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if err == io.EOF {
				break // Terminate the loop when EOF is reached
			}
			sendResponse(conn, 550, "Failed to read file data.")
			return
		}

		if messageType == websocket.BinaryMessage {
			_, err := file.Write(p)
			if err != nil {
				sendResponse(conn, 550, "Failed to write file.")
				return
			}
		}
	}

	// Respond with "226 Transfer complete." when done
	sendResponse(conn, 226, "Transfer complete.")
}
