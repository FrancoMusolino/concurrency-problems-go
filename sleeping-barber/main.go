package main

import (
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var numberOfBarbers = 3
var waitingRoomSize = 4
var shopsOpeningTime = 20 * time.Second

var names = []string{
	"Juan", "Pedro", "Lucía", "María", "Carlos",
	"Sofía", "Andrés", "Valentina", "Miguel", "Camila",
	"Diego", "Paula", "Mateo", "Ana", "Javier",
	"Gabriela", "Fernando", "Laura", "Martín", "Isabella",
	"Emilia", "Tomás", "Julieta", "Agustín", "Victoria",
	"Lucas", "Martina", "Benjamín", "Antonia", "Santiago",
	"Josefina", "Nicolás", "Renata", "Facundo", "Florencia",
	"Manuel", "Mía", "Sebastián", "Lola", "Ramiro",
	"Elena", "Franco", "Bianca", "Simón", "Catalina",
	"Bruno", "Malena", "Iván", "Ariana", "Leandro",
}

// A barber is waiting for a client to show up, meanwhile is sleeping
// when a client shows up, if there is no seats available, he/she leaves
// if there is a seat available and the barber is sleeping, the client wakes up the barber
// if the barber is busy, the client waits on the waiting room
// Once the shop closes, no more clients ara allowed, but the barber has to stay until everyone gets a hair cut

func main() {
	shop := NewShop()
	color.Cyan("%s is open!!", shop.name)

	shopIsClosingChan := make(chan bool)
	shopClosedChan := make(chan bool)

	go func() {
		<-time.After(shopsOpeningTime)
		color.Cyan("Closing %s, no more client are allowed", shop.name)
		shopIsClosingChan <- true
		shop.Close()
		shopClosedChan <- true
	}()

	for i := 0; i < numberOfBarbers; i++ {
		shop.AddBarber(pickRandomName())
	}

	go func() {
		for {
			delay := rand.Intn(500) + 100

			select {
			case <-shopIsClosingChan:
				return
			case <-time.After(time.Duration(delay) * time.Millisecond):
				shop.NewClient(pickRandomName())
			}

		}
	}()

	<-shopClosedChan
	color.Cyan("%s is closed", shop.name)
}

func pickRandomName() string {
	index := rand.Intn(len(names)) + 1
	return names[index-1]
}
