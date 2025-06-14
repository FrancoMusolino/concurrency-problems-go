package main

import (
	"math/rand"
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	name                  string
	barbers               []Barber
	isOpen                bool
	clientsChan           chan Client
	barberDoneWorkingChan chan bool
}

func NewShop() BarberShop {
	return BarberShop{
		name:                  "The Super Barber Shop",
		isOpen:                true,
		clientsChan:           make(chan Client, waitingRoomSize),
		barberDoneWorkingChan: make(chan bool),
	}
}

func (s *BarberShop) AddBarber(name string) {
	barber := Barber{name: name}
	s.barbers = append(s.barbers, barber)
	go barberWorking(s, &barber)
}

func (s *BarberShop) NewClient(name string) {
	c := Client{name: name}

	if !s.isOpen {
		return
	}

	select {
	case s.clientsChan <- c:
		color.Green("\tNew Client %s", c.name)
	default:
		color.Red("\tWaiting room is full, client %s leaves", c.name)
	}

}

func (s *BarberShop) Close() {
	s.isOpen = false
	close(s.clientsChan)

	for i := 0; i < len(s.barbers); i++ {
		<-s.barberDoneWorkingChan
	}

	close(s.barberDoneWorkingChan)
}

type Barber struct {
	name string
}

// Client
type Client struct {
	name string
}

func barberWorking(shop *BarberShop, barber *Barber) {
	var sleepingAt time.Time
	isSleeping := false
	for {
		if len(shop.clientsChan) == 0 {
			isSleeping = true
			sleepingAt = time.Now()
			color.White("\t\tBarber %s sleeping", barber.name)
		}

		c, ok := <-shop.clientsChan
		if !ok {
			color.Magenta("\t\tBarber %s finish working for today", barber.name)
			shop.barberDoneWorkingChan <- true
			return
		}

		if isSleeping {
			isSleeping = false
			elapsed := time.Since(sleepingAt)

			color.White("\t\tAwaking barber %s after sleeping for %.2f seconds", barber.name, elapsed.Seconds())
		}

		delay := rand.Intn(3) + 1
		color.White("\t\tBarber %s doing haircut to client %s. Delay of %d seconds", barber.name, c.name, delay)
		time.Sleep(time.Duration(delay) * time.Second)
	}
}
