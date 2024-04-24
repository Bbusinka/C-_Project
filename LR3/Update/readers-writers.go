package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	M_READERS = 5 
	N_WRITERS = 3 
)

var (
	readersPresenceList []int
	writerPresent       bool
	readersMutex        sync.Mutex
	writersMutex        sync.Mutex
	readersCount        int
)

func sleepRandomTime() {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

func enterLibrary(id int, isReader bool, done chan struct{}) {
	sleepRandomTime()

	if isReader {
		readersMutex.Lock()
		readersCount++
		if readersCount == 1 {
			writersMutex.Lock()
		}
		readersPresenceList = append(readersPresenceList, id)
		fmt.Printf("Czytelnik %d wszedł do czytelni. Obecni czytelnicy: %v. Obecni pisarze: %v\n", id, readersPresenceList, writerPresent)
		readersMutex.Unlock()
	} else {
		writersMutex.Lock()
		for readersCount > 0 {
			writersMutex.Unlock()
			sleepRandomTime()
			writersMutex.Lock()
		}
		writerPresent = true
		fmt.Printf("Pisarz %d wszedł do czytelni. Obecni czytelnicy: %v. Obecni pisarze: %v\n", id, readersPresenceList, writerPresent)
	}

	sleepRandomTime()
}

func exitLibrary(id int, isReader bool, done chan struct{}) {
	if isReader {
		readersMutex.Lock()
		readersCount--
		if readersCount == 0 {
			writersMutex.Unlock()
		}
		readersPresenceList = removeID(readersPresenceList, id)
		fmt.Printf("Czytelnik %d wyszedł z czytelni. Obecni czytelnicy: %v. Obecni pisarze: %v\n", id, readersPresenceList, writerPresent)
		readersMutex.Unlock()
	} else {
		writerPresent = false
		fmt.Printf("Pisarz %d wyszedł z czytelni. Obecni czytelnicy: %v. Obecni pisarze: %v\n", id, readersPresenceList, writerPresent)
		writersMutex.Unlock()
	}

	sleepRandomTime()
}

func removeID(slice []int, id int) []int {
	index := -1
	for i, val := range slice {
		if val == id {
			index = i
			break
		}
	}

	if index != -1 {
		return append(slice[:index], slice[index+1:]...)
	}
	return slice
}

func readerJob(id int, done chan struct{}) {
	for {
		enterLibrary(id, true, done)
		exitLibrary(id, true, done)
	}
}

func writerJob(id int, done chan struct{}) {
	for {
		enterLibrary(id, false, done)
		exitLibrary(id, false, done)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	done := make(chan struct{})

	for i := 0; i < M_READERS; i++ {
		go readerJob(i, done)
	}

	for i := 0; i < N_WRITERS; i++ {
		go writerJob(i, done)
	}

	<-done
}
