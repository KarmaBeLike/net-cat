package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
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
		go handleConn(conn)
	}
}

type client chan<- string // Канал исходящих сообщений
var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // Все входящие сообщения клиента
)

func broadcaster() {
	clients := make(map[client]bool) // Все подключенные клиенты
	for {
		select {
		case msg := <-messages: // Широковещательное входящее сообщение во все каналы исходящих сообщений для клиентов
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)

		}
	}
}

func handleConn(conn net.Conn) {
	Welcom(conn)
	ch := make(chan string) // Исходящие сообщения клиентов
	go clientWriter(conn, ch)
	who := conn.RemoteAddr().String()
	ch <- "\n[ENTER YOUR NAME]: " + who
	messages <- who + " connected"
	entering <- ch
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	leaving <- ch
	messages <- who + " disconnected"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func Welcom(conn net.Conn) {
	file, err := os.ReadFile("logo.txt")
	if err != nil {
		log.Print("File is not correct\n")
		return
	}
	welcomtext := string(file)
	fmt.Fprintf(conn, welcomtext)
}
