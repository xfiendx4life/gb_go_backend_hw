package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/v-lozhkin/backendOneLessons/lesson2/math/mathmaker"
)

type scoreTable struct {
	score map[string]int
	mu    sync.Mutex
}

type client chan<- string

func generateEquation() (mm mathmaker.Mathmaker, eq string) {
	ops := mathmaker.CreateDefaultOperations()
	mm = mathmaker.NewMathmaker()
	eq = mm.MakeMathEquation(ops)
	return
}

func solve(ctx context.Context, attempt chan string, res chan bool, mm chan mathmaker.Mathmaker) {
	// var m mathmaker.Mathmaker
	// подумать, как получить эту штуку и обработать
	for {
		select {
		case <-ctx.Done():
			return
		// case at := <-attempt:
		// 	res <- at == mm.GetEquationResult()
		}
	}
}

func broadcaster(messages chan string, leaving chan client, entering chan client) {
	clients := make(map[client]bool)

	for {
		select {
		case msg := <-messages:
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

func handleConn(conn net.Conn, messages chan string, leaving chan client, entering chan client) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	// who := conn.RemoteAddr().String()
	ch <- "Who are you?"
	input := bufio.NewScanner(conn)
	var who string
	if input.Scan() {
		who = input.Text()
	}
	messages <- who + " has arrived"
	entering <- ch

	// input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	leaving <- ch
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatal(err)
	}
	entering := make(chan client)
	leaving := make(chan client)
	messages := make(chan string)

	go broadcaster(messages, leaving, entering)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, messages, leaving, entering)
	}
}
