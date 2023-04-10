package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type server struct {
	listener       net.Listener
	clients        map[string]net.Conn
	maxConnections int
	connections    int
	mutex          sync.Mutex
}

func NewServer(protocol, addres string, maxConnections int) (*server, error) {
	listener, err := net.Listen(protocol, addres)
	if err != nil {
		return nil, err
	}
	return &server{
		listener:       listener,
		clients:        make(map[string]net.Conn),
		maxConnections: maxConnections,
		mutex:          sync.Mutex{},
	}, nil
}

func main() {
	s, err := NewServer("tcp", ":8000", 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listening on the port: " + "localhost:8000")

	go s.broadcaster()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go s.handleConn(conn)

	}
}

type message struct {
	time string
	name string
	text string
}
type client chan<- string // Канал исходящих сообщений
var (
	clients  = make(map[string]net.Conn) // Все подключенные клиенты
	entering = make(chan message)
	leaving  = make(chan message)
	messages = make(chan message) // Все входящие сообщения клиента
)

const timeFormat = "2006-01-02 15:04:05"

func (s *server) broadcaster() {
	for {
		select {
		case msg := <-messages: // Широковещательное входящее сообщение во все каналы исходящих сообщений для клиентов
			for cli, conn := range s.clients {
				if cli != msg.name {
					conn.Write([]byte(fmt.Sprintf("\n[%s][%s]:%s", msg.time, msg.name, msg.text)))
					conn.Write([]byte(fmt.Sprintf("\n[%s][%s]:", time.Now().Format(timeFormat), cli)))
				}
			}
		case msg := <-entering:
			for cli, conn := range s.clients {
				if cli != msg.name {
					conn.Write([]byte(fmt.Sprintf("\n%s", msg.text)))
					conn.Write([]byte(fmt.Sprintf("\n[%s][%s]:", time.Now().Format(timeFormat), cli)))
				}
			}
		case msg := <-leaving:
			for cli, conn := range s.clients {
				if cli != msg.name {
					conn.Write([]byte(fmt.Sprintf("\n%s", msg.text)))
					conn.Write([]byte(fmt.Sprintf("\n[%s][%s]:", time.Now().Format(timeFormat), cli)))
				}
			}
		}
	}
}

func (s *server) handleConn(conn net.Conn) {
	if s.connections > s.maxConnections {
		conn.Write([]byte("Chatroom is full, try again later"))
		conn.Close()
		return
	}
	s.connections++
	Welcome(conn)
	// ch := make(chan string) // Исходящие сообщения клиентов
	// go clientWriter(conn, ch)
	name := conn.RemoteAddr().String()
	scan := bufio.NewScanner(conn)
	conn.Write([]byte("\n[ENTER YOUR NAME]: "))
	for scan.Scan() {
		name = scan.Text()
		name = strings.TrimSpace(name)
		if len(name) == 0 {
			conn.Write([]byte("incorrect input\n"))
			conn.Write([]byte("[ENTER YOUR NAME]: "))
		} else if _, exist := s.clients[name]; exist {
			conn.Write([]byte("Username is used"))
			conn.Write([]byte("[ENTER YOUR NAME]: "))
		} else {
			break
		}
	}
	s.clients[name] = conn
	entering <- message{
		time: "",
		name: name,
		text: name + " has joined our chat...",
	}
	conn.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format(timeFormat), name)))
	// messages <- name + " connected"
	input := bufio.NewScanner(conn)
	for input.Scan() {
		if input.Text() != "" {
			messages <- message{
				time: time.Now().Format(timeFormat),
				name: name,
				text: input.Text(),
			}
		}
		conn.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format(timeFormat), name)))
	}
	leaving <- message{
		time: "",
		name: name,
		text: name + " has left our chat...",
	}
	// messages <- name + " disconnected"
	conn.Close()
	s.connections--
	s.clients[name].Close()
	delete(s.clients, name)
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func Welcome(conn net.Conn) {
	file, err := os.ReadFile("logo.txt")
	if err != nil {
		log.Print("File is not correct\n")
		return
	}
	welcometext := string(file)
	fmt.Fprintf(conn, welcometext)
}
