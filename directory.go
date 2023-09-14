package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

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

	listing, _ := GenerateDirectoryListing(Dir(ftpContext))

	sendData(conn, listing)

	sendResponse(conn, 226, "Directory listing completed.")
}

func Dir(ftpContext *FTPContext) string {
	str, err := filepath.Abs((fmt.Sprintf("./%s", ftpContext.WorkingDir)))
	if err != nil {
		return ""
	}
	return str
}

func GenerateDirectoryListing(directoryPath string) (string, error) {
	entries, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return "", err
	}

	var listing strings.Builder

	for _, entry := range entries {
		// File type and permissions
		fileInfo, _ := os.Stat(filepath.Join(directoryPath, entry.Name()))
		fileType := "-"
		if fileInfo.IsDir() {
			fileType = "d"
		}

		// Owner, group, and file size
		fileSize := entry.Size()

		// Modification time
		modTime := entry.ModTime().Format("Jan _2 15:04")
		fileName := entry.Name()

		// Append the entry to the listing
		listing.WriteString(fmt.Sprintf("%s %s %s %d %s %s\n", fileType, "-", "-", fileSize, modTime, fileName))
	}

	return listing.String(), nil
}
