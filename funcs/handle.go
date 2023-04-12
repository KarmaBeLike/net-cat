package funcs

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func (s *server) HandleConn(conn net.Conn) {
	s.mutex.Lock()
	flag := s.connections == s.maxConnections
	s.mutex.Unlock()
	if flag {
		conn.Write([]byte("Chatroom is full, try again later"))
		conn.Close()
		return
	}
	s.mutex.Lock()
	s.connections++
	s.mutex.Unlock()

	Welcome(conn)
	var name string
	scan := bufio.NewScanner(conn)
	conn.Write([]byte("\n[ENTER YOUR NAME]: "))
	for scan.Scan() {
		name = scan.Text()
		name = strings.TrimSpace(name)
		if !IsValidName(name) {
			conn.Write([]byte("incorrect\n"))
			conn.Write([]byte("\n[ENTER YOUR NAME]: "))
			continue
		}
		if len(name) == 0 || len(name) > 20 {
			conn.Write([]byte("incorrect input\n"))
			conn.Write([]byte("\n[ENTER YOUR NAME]: "))
		} else if _, exist := s.clients[name]; exist {
			conn.Write([]byte("Username is used\n"))
			conn.Write([]byte("\n[ENTER YOUR NAME]: "))
		} else {
			break
		}
	}
	s.mutex.Lock()
	s.clients[name] = conn
	s.mutex.Unlock()
	s.mutex.Lock()
	entering <- message{
		time: "",
		name: name,
		text: name + " has joined our chat...",
	}

	s.mutex.Unlock()
	input := bufio.NewScanner(conn)
	for input.Scan() {
		if strings.TrimSpace(input.Text()) == "" || len(input.Text()) > 300 {
			conn.Write([]byte("incorrect input\n"))
			conn.Write([]byte((fmt.Sprintf("[%s][%s]:", time.Now().Format(timeFormat), name))))
			continue
		}
		messages <- message{
			time: time.Now().Format(timeFormat),
			name: name,
			text: strings.TrimSpace(input.Text()),
		}
		text := fmt.Sprintf("[%s][%s]:%s\n", time.Now().Format(timeFormat), name, strings.TrimSpace(input.Text()))
		conn.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format(timeFormat), name)))
		s.mutex.Lock()
		history = append(history, text)
		s.mutex.Unlock()
	}
	leaving <- message{
		time: "",
		name: name,
		text: name + " has left our chat...",
	}
	// messages <- name + " disconnected"
	conn.Close()
	s.mutex.Lock()
	s.connections--
	s.clients[name].Close()
	delete(s.clients, name)
	s.mutex.Unlock()
}

func IsValidName(name string) bool {
	for _, c := range name {
		if c < 32 || c > 126 {
			return false
		}
	}
	return true
}
