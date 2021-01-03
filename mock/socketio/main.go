package main

import (
	"bufio"
	socketio_client "github.com/zhouhui8915/go-socket.io-client"
	"log"
	"os"
)



func main() {
	opts := &socketio_client.Options{
		//Transport:"polling",
		Transport:"websocket",
		Query:     make(map[string]string),
	}
	opts.Query["user"] = "user"
	opts.Query["pwd"] = "pass"
	uri := "http://127.0.0.1:9090"

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		log.Printf("NewClient error:%v\n", err)
		return
	}

	client.On("error", func() {
		log.Printf("on error\n")
	})
	client.On("connection", func() {
		log.Printf("on connect\n")
	})

	client.On("minerInfo", func(msg string) {
		log.Printf("on minerInfo:%v\n", msg)
	})

	client.On("subMinerInfo", func(msg string) {
		log.Printf("on subMinerInfo:%v\n", msg)
	})

	client.On("disconnection", func() {
		log.Printf("on disconnect\n")
	})

	reader := bufio.NewReader(os.Stdin)
	for {
		data, _, _ := reader.ReadLine()
		command := string(data)
		client.Emit("minerInfo", command)
		client.Emit("subMinerInfo", command)
		log.Printf("send message:%v\n", command)
	}
}