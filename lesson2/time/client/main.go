package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8001")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//buf := make([]byte, 256)
	for {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			log.Println(err)
		}
	}
}
