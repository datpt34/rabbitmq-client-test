package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// Thông tin kết nối
	host := flag.String("h", "171.251.89.96", "Địa chỉ host")
	port := flag.String("p", "5672", "Cổng")
	username := flag.String("u", "abc", "Tên người dùng")
	password := flag.String("password", "TZ8EveBPR9GxkEb89ayPjRrfLA0LTvBq", "Mật khẩu")

	// Tạo chuỗi kết nối
	connectionString := "amqp://" + *username + ":" + *password + "@" + *host + ":" + *port + "/"
	fmt.Println(connectionString)
	// Tạo kết nối đến RabbitMQ server
	conn, err := amqp.Dial(connectionString)
	failOnError(err, "Failed to connect to RabbitMQ")

	ticker := time.NewTicker(1 * time.Second) // Tạo một đối tượng Ticker với khoảng thời gian 1 giây
	defer ticker.Stop()                       // Đảm bảo dừng Ticker khi chương trình kết thúc

	defer conn.Close()
	log.Println("Kết nối đến RabbitMQ thành công")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	//q, err := ch.QueueDeclare(
	//	"response", // name
	//	false,      // durable
	//	false,      // delete when unused
	//	false,      // exclusive
	//	false,      // no-wait
	//	nil,        // arguments
	//)
	failOnError(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		"response", // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	//go func() {
	//	start := time.Now()
	//	for {
	//		select {
	//		case <-done: // Nếu nhận được tín hiệu kết thúc từ channel done, thoát khỏi vòng lặp
	//			return
	//		case <-ticker.C: // Khi Ticker gửi tín hiệu qua channel ticker.C, in ra màn hình
	//			fmt.Println("Time elapsed: ", time.Since(start))
	//		default:
	//			if len(msgs) > 0 {
	//				for d := range msgs {
	//					log.Printf("Received a message: %s", d.Body)
	//				}
	//			}
	//		}
	//	}
	//}()
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	// Thực hiện các thao tác với kết nối đã tạo ở đây
	// Ví dụ: tạo kênh, gửi hoặc nhận tin nhắn, ...
}
