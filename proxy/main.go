package main

import (
	"bytes"
	"fmt"
	"github.com/valinurovam/garagemq/amqp"
	"log"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		client, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		rabbitmq, err := net.Dial("tcp", "localhost:5672")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequestClientToRabbitMq(client, rabbitmq)
		go handleRequestRabbitMqToClient(rabbitmq, client)
	}
}

func handleRequestClientToRabbitMq(r net.Conn, w net.Conn) {
	defer r.Close()
	defer w.Close()
	buf := make([]byte, 8)
	_, err := r.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("first read: ", string(buf))
	w.Write(buf)
	for {
		// incoming request
		buffer := make([]byte, 1024)
		n, err := r.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		data := buffer[:n]
		fmt.Printf("read from client %d bytes, data: %s \n", n, string(data))
		// write data to response
		newReader := bytes.NewReader(data)
		f, err := amqp.ReadFrame(newReader)
		if err != nil {
			log.Fatal("error ReadFrame", err)
		}
		err = authorize(f)
		if err != nil {
			log.Fatal("error authorize")
		}
		_, err = w.Write(data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("write to RabbitMq ok")

		// close conn
	}

}

func authorize(f *amqp.Frame) error {
	buffer := bytes.NewReader([]byte{})
	if f != nil {
		fmt.Println("frame Payload: ", string(f.Payload))
		switch f.Type {
		case amqp.FrameMethod:
			fmt.Println("Frame Method")
			buffer.Reset(f.Payload)
			method, err := amqp.ReadMethod(buffer, "amqp-rabbit")
			fmt.Println("Incoming method <- " + method.Name())
			if err != nil {
				fmt.Println("Error on handling frame")
				return err
			}

			if err := handleMethod(method); err != nil {
				return err
			}
		case amqp.FrameHeader:
			fmt.Println("Frame Header")
			return nil
		case amqp.FrameBody:
			fmt.Println("FrameBody")
			return nil

		}
	}
	return nil
}

func handleMethod(method amqp.Method) error {
	fmt.Println("handle Method")
	switch method.ClassIdentifier() {
	case amqp.ClassConnection:
		fmt.Println("ClassConnection")
		return nil
	case amqp.ClassChannel:
		fmt.Println("ClassChannel")
		return nil
	case amqp.ClassBasic:
		fmt.Println("ClassBasic")
		return nil
	case amqp.ClassExchange:
		fmt.Println("ClassExchange")
		return nil
	case amqp.ClassQueue:
		fmt.Println("ClassQueue")
		return nil
	case amqp.ClassConfirm:
		fmt.Println("ClassConfirm")
		return nil
	}
	return nil
}
func handleRequestRabbitMqToClient(r net.Conn, w net.Conn) {
	defer r.Close()
	defer w.Close()
	for {
		// incoming request
		buffer := make([]byte, 1024)
		n, err := r.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		data := buffer[:n]
		fmt.Printf("read from RabbitMq %d bytes, data: %s \n", n, string(data))

		// write data to response
		_, err = w.Write(data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("write to client ok")

		// close conn
	}

}
