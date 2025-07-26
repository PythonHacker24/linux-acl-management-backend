// wsclient.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run wsclient.go <url> <jwt_token>")
		return
	}

	url := os.Args[1]
	jwt := os.Args[2]

	header := http.Header{}
	header.Add("Authorization", "Bearer "+jwt)

	urlQ := fmt.Sprintf("%s?token=%s", url, jwt)

	conn, resp, err := websocket.DefaultDialer.Dial(urlQ, header)
	if err != nil {
		log.Fatalf("Dial failed: %v, response: %+v", err, resp)
	}
	defer conn.Close()

	fmt.Println("Connected to", url)
	fmt.Println("Type a message and press Enter to send. Ctrl+C to quit.")

	// Read incoming messages in a goroutine
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				os.Exit(1)
			}
			fmt.Println("Received:", string(msg))
		}
	}()

	// Read user input and send messages
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		err := conn.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}
