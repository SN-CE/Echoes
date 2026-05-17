package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"log"
	"time"
	"net/http"
	"github.com/gorilla/websocket"
	"sync"
)

var client_list map[string]*websocket.Conn
var mutex_key sync.RWMutex

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func choose_username(conn *websocket.Conn) (string, error) {
	err := conn.WriteMessage(websocket.TextMessage, []byte("Hey! Before you start chatting, I need you to register a username.\n\nThis so your friends know who you are, instead of seeing your IP Address."))
	if err != nil {
		return "", err
	}

	mutex_key.RLock()
	var user_count int = 0
	for client, _ := range client_list {
		err = conn.WriteMessage(websocket.TextMessage, []byte("[yellow::b]" + client))
		if err != nil {
			mutex_key.RUnlock()
			return "", err
		}
		user_count++
	}
	mutex_key.RUnlock()

	if user_count == 0 {
		err = conn.WriteMessage(websocket.TextMessage, []byte("[yellow::b]No other users Connected."))
		if err != nil {
			return "", err
		}
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte("[white]Please type your unique username."))
	if err != nil {
		return "", err
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return "", err
		}
		username := string(msg)
		mutex_key.RLock()
		_, exists := client_list[username]
		mutex_key.RUnlock()
		if exists {
			err := conn.WriteMessage(websocket.TextMessage, []byte("[red]A client with that name already exists. Please choose another username :"))
			if err != nil {
				return "", err
			}
			continue
		}
		err = conn.WriteMessage(websocket.TextMessage, []byte("[red::b]\n\nWELCOME TO ECHOES[-::-]"))
		if err != nil {
			return "", err
		}
		return username, nil
	}
}

func remove_client(client string) {
	mutex_key.Lock()
	delete(client_list, client)
	mutex_key.Unlock()
}

func chat_handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Could Not Upgrade Connection. Error : %v\n", err)
		return
	}

	defer conn.Close()

	client_ip := conn.RemoteAddr().String()
	client, err := choose_username(conn)

	if err != nil {
		fmt.Printf("Client %s disconnected.\n", client_ip)
		return
	}

	mutex_key.Lock()
	client_list[client] = conn
	mutex_key.Unlock()

	fmt.Printf("Client %s connected as %s.\n", client_ip, client)

	var disconn_clients []string

	mutex_key.RLock()
	for client_copy, conn_copy := range client_list {
		if client_copy == client {
			continue
		} else {
			err := conn_copy.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[blue::b]Client %s connected.", client)))
			if err != nil {
				disconn_clients = append(disconn_clients, client_copy)
				fmt.Printf("Client %s disconnected.\n", client_copy)
			}
		}
	}
	mutex_key.RUnlock()

	for _, d_client := range disconn_clients {
			remove_client(d_client)
	}

	mutex_key.RLock()
	for _, d_client := range disconn_clients {
		for _, conn_copy := range client_list {
			conn_copy.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[blue::b]Client %s disconnected.", d_client)))
		}
	}
	mutex_key.RUnlock()

	disconn_clients = nil

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Client %s disconnected.\n", client)

			mutex_key.RLock()
			for client_copy, conn_copy := range client_list {
				if client_copy == client {
					continue
				} else {
					conn_copy.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[blue::b]Client %s disconnected.", client)))
				}
			}
			mutex_key.RUnlock()

			remove_client(client)
			return
		}

		msg_text := client + " : " + string(msg)

		mutex_key.RLock()
		for client_copy, conn_copy := range client_list {
			if client_copy == client {
				err := conn_copy.WriteMessage(websocket.TextMessage, []byte("[green::b]You : " + string(msg)))
				if err != nil {
					disconn_clients = append(disconn_clients, client_copy)
					fmt.Printf("Client %s disconnected.\n", client_copy)
				}
			} else {
				err := conn_copy.WriteMessage(websocket.TextMessage, []byte("[white]" + msg_text))
				if err != nil {
					disconn_clients = append(disconn_clients, client_copy)
					fmt.Printf("Client %s disconnected.\n", client_copy)
				}
			}
		}
		mutex_key.RUnlock()
		for _, d_client := range disconn_clients {
			remove_client(d_client)
		}

		mutex_key.RLock()
		for _, d_client := range disconn_clients {
			for _, conn_copy := range client_list {
				conn_copy.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[blue::b]Client %s disconnected.", d_client)))
			}
		}
		mutex_key.RUnlock()

		disconn_clients = nil

		fmt.Printf("Message From %s : %s\n", client, string(msg))
	}
}

func main() {
	client_list = make(map[string]*websocket.Conn)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.HandleFunc("/", chat_handler)

	var port_no string
	fmt.Printf("Please Enter the Port you wish to run the Echoes Server at - ")
	fmt.Scanln(&port_no)

	server := &http.Server {
		Addr: ":" + port_no,
		Handler: mux,
	}

	fmt.Printf("Echoes Server Started on Port %s.\n", port_no)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Fatal Error. Server Stopped.\n")
			return
		}
	}()

	<-stop

	mutex_key.Lock()
	fmt.Printf("Caught Signal. Exiting...\n")
	for client, conn := range client_list {
		conn.Close()
		delete(client_list, client)
	}
	mutex_key.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
