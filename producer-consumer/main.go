package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAX_NOTIFICATIONS = 5

var sendedNotifications int

type Notification struct {
	id      int
	userId  int
	message string
}

type Producer struct {
	data chan Notification
	quit chan struct{}
}

func (p *Producer) Close() {
	p.quit <- struct{}{}
}

type User struct {
	id   int
	name string
}

func produceNotification(notificationId, userId int) *Notification {
	if sendedNotifications >= MAX_NOTIFICATIONS {
		return &Notification{}
	}

	delay := rand.Intn(3) + 1
	fmt.Printf("Produciendo notificación para el usuario %d, tomará %d segundos...\n", userId, delay)
	time.Sleep(time.Duration(delay) * time.Second)

	return &Notification{
		id:      notificationId,
		userId:  userId,
		message: fmt.Sprintf("Notificación para el usuario %d", userId),
	}
}

func pickRandomUser(users []User) int {
	return rand.Intn(len(users)) + 1
}

func listenForNotifications(producer *Producer, users []User) {
	i := 0

	for {
		i++
		notification := produceNotification(i, pickRandomUser(users))

		select {
		case producer.data <- *notification:
		case <-producer.quit:
			close(producer.data)
			close(producer.quit)
			return
		}

	}
}

func consume(producer *Producer) {
	for notification := range producer.data {
		if notification.id >= MAX_NOTIFICATIONS {
			producer.Close()
			fmt.Println("Closing channel")
			return
		}

		fmt.Printf("Se recibió la notificación N° %d del usuario %d: %s\n", notification.id, notification.userId, notification.message)
	}
}

func main() {
	var users []User
	names := []string{"Juan", "Pepe", "Sofía", "Luciano"}
	for i, name := range names {
		users = append(users, User{
			id:   i + 1,
			name: name,
		})
	}

	producer := Producer{
		data: make(chan Notification),
		quit: make(chan struct{}),
	}

	go listenForNotifications(&producer, users)
	consume(&producer)
}
