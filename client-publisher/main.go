package main

import (
	"context"
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
	//username := "admin"
	//password := "admin"
	//host := "localhost" // Thay đổi thành địa chỉ RabbitMQ server thực tế
	//port := "8080"      // Thay đổi nếu cổng mặc định đã được sửa đổi

	host := flag.String("h", "117.1.29.217", "Địa chỉ host")
	port := flag.String("p", "5672", "Cổng")
	username := flag.String("u", "abc", "Tên người dùng")
	password := flag.String("password", "46zGKA5FSEHdfWhVHM33ViPLVRMMwIUc", "Mật khẩu")

	// Tạo chuỗi kết nối
	connectionString := "amqp://" + *username + ":" + *password + "@" + *host + ":" + *port + "/"
	fmt.Println("Connecting: ", connectionString)
	// Tạo kết nối đến RabbitMQ server
	conn, err := amqp.Dial(connectionString)
	failOnError(err, "Failed to connect to RabbitMQ")

	// Đảm bảo dừng Ticker khi chương trình kết thúc

	done := make(chan bool) // Tạo một channel để thông báo khi kết thúc
	defer conn.Close()
	log.Println("Kết nối đến RabbitMQ thành công")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var forever chan struct{}

	go startPublisher(ctx, done, ch)

	log.Printf(" [*] Publishing messages. To exit press CTRL+C")
	<-forever
}

func startPublisher(ctx context.Context, done chan bool, ch *amqp.Channel) {
	start := time.Now()
	ticker := time.NewTicker(1 * time.Second) // Tạo một đối tượng Ticker với khoảng thời gian 1 giây
	defer ticker.Stop()
	for {
		select {
		case <-done: // Nếu nhận được tín hiệu kết thúc từ channel done, thoát khỏi vòng lặp
			return
		case <-ticker.C: // Khi Ticker gửi tín hiệu qua channel ticker.C, in ra màn hình
			fmt.Println("Time elapsed: ", time.Since(start))
			body := fmt.Sprintf("Hello World! %v", time.Since(start))
			err := ch.PublishWithContext(ctx,
				"exchange_1", // exchange
				"queue_1",    // routing key
				false,        // mandatory
				false,        // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			failOnError(err, "Failed to publish a message")
			log.Printf(" [x] Sent %s\n", body)
		}
	}
}
