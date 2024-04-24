package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

const (
    m = 10 // Liczba rzędów
    n = 10 // Liczba kolumn
    k = m * n - 1 // Maksymalna liczba podróżników
)

type Traveler struct {
    id int
    x  int
    y  int
    possibleMoves [][]int
}

func main() {

    rand.Seed(time.Now().UnixNano())

    kratowaPlansza := make([][]int, m)
    for i := range kratowaPlansza {
        kratowaPlansza[i] = make([]int, n)
    }

    var mutex sync.Mutex // Mutex do synchronizacji dostępu do planszy

    var wg sync.WaitGroup // WaitGroup do śledzenia aktywnych goroutines

    // Inicjalilzacja listy podróżników
    travelers := make([]*Traveler, 0)

    // Goroutine, aby zgenerować podróżników
    go func() {
        for currentID := 1; currentID <= k; currentID++ {
            x, y := rand.Intn(m), rand.Intn(n)

            mutex.Lock()
            for kratowaPlansza[x][y] != 0 {
                x, y = rand.Intn(m), rand.Intn(n)
            }
            kratowaPlansza[x][y] = currentID
            mutex.Unlock() 

            newTraveler := &Traveler{id: currentID, x: x, y: y}
            newTraveler.possibleMoves = getEmptyNeighbors(x, y, kratowaPlansza)
            travelers = append(travelers, newTraveler)
           
            // Goroutine dla każdego podróżnika
            wg.Add(1)
            go func(traveler *Traveler) {
                defer wg.Done()
                simulateTraveler(traveler, kratowaPlansza, &mutex)
            }(newTraveler)

            time.Sleep(time.Second)

           

        }
    }()

       // Goroutine, aby wyświetlić tablicę
    go func() {
        for {
            mutex.Lock() 
            printKrate(kratowaPlansza, travelers)
            mutex.Unlock() 
            time.Sleep(2 * time.Second) 
        }
    }()
    
    for {
        time.Sleep(2 * time.Second)
    }

    // Oczekiwanie na zakończenie wszystkich procedur podróżnych
    go func() {
        wg.Wait()
        
        for {
            time.Sleep(2 * time.Second)
        }
    }()

    wg.Wait()
}

func simulateTraveler(traveler *Traveler, kratowaPlansza [][]int, mutex *sync.Mutex) {
    for {
        time.Sleep(time.Second) 

        mutex.Lock()

        // Symulacja zmiany pozycji podróżnika
        if len(traveler.possibleMoves) > 0 {
           // Wybieramy losową pustą przestrzeń z sąsiednich komórek
            moveIndex := rand.Intn(len(traveler.possibleMoves))
            newX, newY := traveler.possibleMoves[moveIndex][0], traveler.possibleMoves[moveIndex][1]

            kratowaPlansza[traveler.x][traveler.y] = 0
            kratowaPlansza[newX][newY] = traveler.id

            // Aktualizacja pozycji podróżnika
            traveler.x = newX
            traveler.y = newY

            // Czyszczenie possibleMoves
            traveler.possibleMoves = [][]int{}
        }

        mutex.Unlock() 
    }
}

func getEmptyNeighbors(x, y int, kratowaPlansza [][]int) [][]int {
    neighbors := [][]int{}

    // Lewo
    if x-1 >= 0 && kratowaPlansza[x-1][y] == 0 {
        neighbors = append(neighbors, []int{x - 1, y})
    }
    // Prawo
    if x+1 < m && kratowaPlansza[x+1][y] == 0 {
        neighbors = append(neighbors, []int{x + 1, y})
    }
    // Dół
    if y-1 >= 0 && kratowaPlansza[x][y-1] == 0 {
        neighbors = append(neighbors, []int{x, y - 1})
    }
    // Góra
    if y+1 < n && kratowaPlansza[x][y+1] == 0 {
        neighbors = append(neighbors, []int{x, y + 1})
    }

    return neighbors
}

func printKrate(kratowaPlansza [][]int, travelers []*Traveler) {

    fmt.Println("Stan planszy:")

    for i := 0; i < m; i++ {
        for j := 0; j < n; j++ {
            foundTraveler := false
            for _, traveler := range travelers {
                if traveler.x == i && traveler.y == j {
                    fmt.Printf("%02d  ", traveler.id)
                    foundTraveler = true
                    break
                }
            }
            if !foundTraveler {
                if kratowaPlansza[i][j] == 0 {
                    fmt.Print(" .  ") // Wolne miejsce
                } else {
                    fmt.Printf("%02d  ", kratowaPlansza[i][j])
                }
            }

            // Wyświetlanie granicy
            if j < n-1 {
                fmt.Print("|")
            }
        }
        fmt.Println()
        if i < m-1 {
            for j := 0; j < n; j++ {
                fmt.Print("----")
                if j < n-1 {
                    fmt.Print("+")
                }
            }
            fmt.Println()
        }
    }
}
