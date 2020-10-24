// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const timeout = 5 * time.Minute

type client struct {
	name string
	ch   chan<- string //服务端和客户端的交流通道
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn) //处理连接的每一个客户端
	}
}
func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.ch <- msg
			}
		case cli := <-entering:
			clients[cli] = true
			//记录在线的用户
			var onlines []string
			for c := range clients {
				onlines = append(onlines, c.name)
			}
			cli.ch <- fmt.Sprintf("%d clients: %s", len(clients), strings.Join(onlines, ", "))
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}
func handleConn(conn net.Conn) {
	talk := make(chan struct{}) //监测客户端的活动

	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- client{who, ch} //通知客户端的到来

	input := bufio.NewScanner(conn)
	go func() {
		for {
			if input.Scan() {
				messages <- who + ": " + input.Text()
				talk <- struct{}{}
			} else {
				leaving <- client{who, ch}
				messages <- who + " has left"
				conn.Close()
				return
			}
		}

	}()

	for {
		select {
		case _, ok := <-talk:
			if !ok {
				leaving <- client{who, ch}
				messages <- who + " has left"
				conn.Close()
				return
			}
		case <-time.After(timeout):
			leaving <- client{who, ch}
			messages <- who + " has left"
			conn.Close()
			return
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
