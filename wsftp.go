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
	cmd := strings.Fields(command)
	if len(cmd) == 0 {
		return
	}

	ftpContext := getFTPContext(conn)
	if ftpContext == nil {
		ftpContext = &FTPContext{}
		setFTPContext(conn, ftpContext)
	}

	switch cmd[0] {
	case "USER":
		handleUSER(conn, cmd[1])
	case "PASS":
		handlePASS(conn, cmd[1])
	case "CWD":
		handleCWD(conn, cmd[1])
	case "PWD":
		handlePWD(conn)
	case "LIST":
		handleLIST(conn)
	case "RETR":
		handleRETR(conn, cmd[1])
	case "STOR":
		handleSTOR(conn, cmd[1])
	case "QUIT":
		handleQUIT(conn)
	default:
		sendResponse(conn, 502, "Command not implemented.")
	}
}

func main() {
	http.HandleFunc("/", handleFTPConnection)
	fmt.Println("FTP server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
