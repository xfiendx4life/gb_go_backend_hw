package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/v-lozhkin/backendOneLessons/lesson2/math/mathmaker"
)

type score map[string]int

type client chan<- string

func generateEquation() (mm mathmaker.Mathmaker, eq string) {
	ops := mathmaker.CreateDefaultOperations()
	mm = mathmaker.NewMathmaker()
	eq = mm.MakeMathEquation(ops)
	return
}

func shareToClients(clients map[client]bool, msg string) {
	for cli := range clients {
		cli <- msg
	}
}

func scorePrep(sc score) string {
	table := "Score:\n"
	for k, v := range sc {
		table += fmt.Sprintf("%s: %d\n", k, v)
	}
	return table
}

func broadcaster(messages chan string, leaving chan client, entering chan client) {
	clients := make(map[client]bool)
	mm, eq := generateEquation()
	sc := make(score)
	for {
		select {
		case msg := <-messages:
			shareToClients(clients, msg)
			if strings.HasSuffix(msg, "has arrived") {
				shareToClients(clients, fmt.Sprintf("Solve this:\n %s", eq))
				log.Println(mm.GetEquationResult())
			}
			if strings.Contains(msg, ":") && strings.Split(msg, ": ")[1] == mm.GetEquationResult() {
				mm, eq = generateEquation()
				user := strings.Split(msg, ":")[0]
				sc[user] += 1
				shareToClients(clients, fmt.Sprintf("User %s won this round,\n%s\nthe next one is:\n %s", user, scorePrep(sc), eq))
				log.Println(mm.GetEquationResult())
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
	entering <- ch
	messages <- who + " has arrived"

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
