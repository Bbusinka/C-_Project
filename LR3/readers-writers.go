package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ReaderWriterSystem struct {
	readersPresent   []int
	writersPresent   []int
	readersPresentMu sync.Mutex
	writersPresentMu sync.Mutex
}

func (sys *ReaderWriterSystem) enterReader(id int) {
	sys.readersPresentMu.Lock()
	sys.readersPresent = append(sys.readersPresent, id)
	fmt.Printf("Czytelnik %d wszedł do czytelni. Obecni czytelnicy: %v. Obecni pisarze: %v\n", id, sys.readersPresent, sys.writersPresent)
	sys.readersPresentMu.Unlock()
	time.Sleep(time.Millisecond * time.Duration(randomInt(500, 1000)))
}

func (sys *ReaderWriterSystem) exitReader(id int) {
	sys.readersPresentMu.Lock()
	for i, readerID := range sys.readersPresent {
		if readerID == id {
			sys.readersPresent = append(sys.readersPresent[:i], sys.readersPresent[i+1:]...)
			break
		}
	}
	fmt.Printf("Czytelnik %d wyszedł z czytelni. Obecni czytelnicy: %v. Obecni pisarze: %v\n", id, sys.readersPresent, sys.writersPresent)
	sys.readersPresentMu.Unlock()
	time.Sleep(time.Millisecond * time.Duration(randomInt(500, 1000)))
}

func (sys *ReaderWriterSystem) enterWriter(id int) {
	sys.writersPresentMu.Lock()
	sys.writersPresent = append(sys.writersPresent, id)
	fmt.Printf("Pisarz %d wszedł do czytelni. Obecni czytelnicy: %v. Obecni pisarze: %v\n", id, sys.readersPresent, sys.writersPresent)
	sys.writersPresentMu.Unlock()
	time.Sleep(time.Millisecond * time.Duration(randomInt(500, 1000)))
}

func (sys *ReaderWriterSystem) exitWriter(id int) {
	sys.writersPresentMu.Lock()
	for i, writerID := range sys.writersPresent {
		if writerID == id {
			sys.writersPresent = append(sys.writersPresent[:i], sys.writersPresent[i+1:]...)
			break
		}
	}
	fmt.Printf("Pisarz %d wyszedł z czytelni. Obecni czytelnicy: %v. Obecni pisarze: %v\n", id, sys.readersPresent, sys.writersPresent)
	sys.writersPresentMu.Unlock()
	time.Sleep(time.Millisecond * time.Duration(randomInt(500, 1000)))
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func reader(id int, sys *ReaderWriterSystem, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		time.Sleep(time.Millisecond * time.Duration(randomInt(1000, 2000)))
		sys.enterReader(id)
		sys.exitReader(id)
	}
}

func writer(id int, sys *ReaderWriterSystem, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		time.Sleep(time.Millisecond * time.Duration(randomInt(1000, 2000)))
		sys.enterWriter(id)
		sys.exitWriter(id)
	}
}

func main() {
	sys := ReaderWriterSystem{}
	var wg sync.WaitGroup

	m := 5 // liczba czytelników
	n := 3 // liczba pisarzy

	// Start readers
	for i := 0; i < m; i++ {
		wg.Add(1)
		go reader(i, &sys, &wg)
	}

	// Start writers
	for i := 0; i < n; i++ {
		wg.Add(1)
		go writer(i, &sys, &wg)
	}

	wg.Wait()
}
