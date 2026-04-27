package main

import (
	"bufio"
	"net"
	"os"
	"time"
)

func main() {
	sock := os.Getenv("SOCKET_PATH")
	if sock == "" {
		os.Exit(1)
	}

	conn, err := net.DialTimeout("unix", sock, time.Second)
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(time.Second))

	if _, err := conn.Write([]byte("GET /ready HTTP/1.0\r\nHost: x\r\n\r\n")); err != nil {
		os.Exit(1)
	}

	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil || len(line) < 12 || line[9] != '2' {
		os.Exit(1)
	}
}
