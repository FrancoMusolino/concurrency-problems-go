package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const NUM_OF_PHILOSOPHERS = 5

var orderMu = sync.Mutex{}
var orderFinishEating []string

var hunger = 3

type Fork struct {
	position    int
	isBeingUsed bool
	mu          *sync.Mutex
}

func (f *Fork) use() {
	f.mu.Lock()
	f.isBeingUsed = true
}

func (f *Fork) free() {
	f.isBeingUsed = false
	f.mu.Unlock()
}

type Philosopher struct {
	name      string
	rightFork *Fork
	leftFork  *Fork
}

func eat(philosopher *Philosopher) {
	for i := hunger; i > 0; i-- {
		if philosopher.rightFork.position > philosopher.leftFork.position {
			philosopher.leftFork.use()
			philosopher.rightFork.use()
		} else {
			philosopher.rightFork.use()
			philosopher.leftFork.use()
		}

		delay := rand.Intn(3) + 1

		fmt.Printf("\tPhilosopher %s eating\n", philosopher.name)
		time.Sleep(time.Duration(delay) * time.Second)

		philosopher.leftFork.free()
		philosopher.rightFork.free()

		fmt.Printf("\tPhilosopher %s sleeping\n", philosopher.name)
		time.Sleep(time.Duration(delay) * time.Second)
	}

	fmt.Printf("Philosopher %s is satisfied\n", philosopher.name)

	orderMu.Lock()
	orderFinishEating = append(orderFinishEating, philosopher.name)
	orderMu.Unlock()
}

func main() {
	fmt.Println("Starting dining")
	fmt.Println("---------------")

	var wg sync.WaitGroup
	philosophers := prepareTable()

	for _, philosopher := range philosophers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			eat(&philosopher)
		}()
	}

	wg.Wait()

	fmt.Println("---------------")
	fmt.Println("Ending dining")

	fmt.Println("The philosophers order of finish eating is: ")
	for _, p := range orderFinishEating {
		fmt.Printf("%s, ", p)
	}
}

func prepareTable() []Philosopher {
	var names = []string{"Juan", "Pedro", "Sofía", "Matías", "Juliana"}
	var philosophers []Philosopher
	var forks []Fork

	for i := 0; i < NUM_OF_PHILOSOPHERS; i++ {
		forks = append(forks, Fork{
			position: i,
			mu:       &sync.Mutex{},
		})
	}

	for i := 0; i < NUM_OF_PHILOSOPHERS; i++ {
		left, right := computePhilosopherForks(i)
		philosophers = append(philosophers, Philosopher{
			name:      names[i],
			leftFork:  &forks[left],
			rightFork: &forks[right],
		})
	}

	return philosophers
}

func computePhilosopherForks(position int) (left int, right int) {
	return position, (position + 1) % NUM_OF_PHILOSOPHERS
}
