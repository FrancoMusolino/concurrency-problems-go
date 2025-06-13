package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const NUMBER_OF_BARBERS = 2
const WAITING_ROOM_SIZE = 4

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

// Barber shop
type BarberShop struct {
	name    string
	barbers []Barber
	isOpen  bool
	ch      chan Client
	mu      *sync.Mutex
}

func (bs *BarberShop) Close() {
	bs.mu.Lock()
	bs.isOpen = false
	bs.mu.Unlock()

	close(bs.ch)
}

func (bs *BarberShop) IsClose() bool {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	return !bs.isOpen
}

// Barbers
type Barber struct {
	name string
}

// Client
type Client struct {
	name string
}

func barberWork(b *Barber, ch <-chan Client, wg *sync.WaitGroup) {
	for {
		c, ok := <-ch
		if !ok {
			return
		}

		delay := rand.Intn(5) + 1

		fmt.Printf("\t\tBarber %s doing haircut to client %s. Delay of %d seconds\n", b.name, c.name, delay)
		time.Sleep(time.Duration(delay) * time.Second)
		wg.Done()
	}
}

func main() {
	wg := sync.WaitGroup{}
	barberShop := initBarberShop()
	fmt.Printf("%s is open!!\n", barberShop.name)

	go func() {
		time.Sleep(20 * time.Second)
		fmt.Printf("Closing %s, no more client are allowed\n", barberShop.name)
		barberShop.Close()
	}()

	for _, b := range barberShop.barbers {
		go barberWork(&b, barberShop.ch, &wg)
	}

	for {
		time.Sleep(1 * time.Second)
		if barberShop.IsClose() {
			break
		}

		client := Client{
			name: pickRandomName(),
		}
		fmt.Println("\t New client", client.name)

		wg.Add(1)
		barberShop.ch <- client
	}

	wg.Wait()
}

func initBarberShop() BarberShop {
	var barbers []Barber

	for i := 0; i < NUMBER_OF_BARBERS; i++ {
		barbers = append(barbers, Barber{
			name: pickRandomName(),
		})
	}

	return BarberShop{
		name:    "The Super Barber Shop",
		isOpen:  true,
		barbers: barbers,
		ch:      make(chan Client, WAITING_ROOM_SIZE-1),
		mu:      &sync.Mutex{},
	}
}

func pickRandomName() string {
	index := rand.Intn(len(names)) + 1
	return names[index-1]
}
