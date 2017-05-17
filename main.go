package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	hostname string
	metrics  = &Metrics{
		data:   make(map[string]map[string]interface{}),
		inputs: make(map[string]func() interface{}),
	}

	settings struct {
		listenAddr string
		updateInt  int
	}
)

func init() {
	hostname, _ = os.Hostname()
	metrics.data[hostname] = make(map[string]interface{})

	flag.StringVar(&settings.listenAddr, "listen", "localhost:8080", "Listen address:port")
	flag.IntVar(&settings.updateInt, "update-int", 30, "Metrics update interval")
	flag.Parse()
}

func run() {
	metrics.fetchMetrics()
	ticker := time.NewTicker(time.Second * time.Duration(settings.updateInt))
	go func() {
		for _ = range ticker.C {
			metrics.fetchMetrics()
		}
	}()
}

func listen(bindAddr string) {
	server, err := net.Listen("tcp", bindAddr)
	if err != nil {
		log.Fatalf("Listener error: %s\n", err)
	} else {
		log.Printf("Listing on %s\n", bindAddr)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Printf("Server error: %s\n", err)
			continue
		}
		reqHandler(conn)
	}
}

func reqHandler(conn net.Conn) {
	reqBuf := make([]byte, 8)
	mlen, err := conn.Read(reqBuf)
	if err != nil && err != io.EOF {
		fmt.Println(err.Error())
	}

	req := strings.TrimSpace(string(reqBuf[:mlen]))
	log.Printf("%s command received from %s\n",
		req, strings.Split(conn.RemoteAddr().String(), ":")[0])

	switch req {
	case "stats":
		r, _ := metrics.getMetrics()
		conn.Write(r)
		conn.Close()
	default:
		m := fmt.Sprintf("Not a command: %s\n", req)
		conn.Write([]byte(m))
		conn.Close()
	}
}

func main() {
	run()
	listen(settings.listenAddr)
}
