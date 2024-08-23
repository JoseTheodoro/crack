package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/JoseTheodoro42/carck/internal/service"
)

func worker(workerID int, partCh chan service.Part, wg *sync.WaitGroup, service *service.Csv) {
	defer wg.Done()

	for part := range partCh {

		fileContent, err := service.ReadFile(part)
		fmt.Printf("File read successfully by worker: %d\n", workerID)

		service.WriteFilePart(fileContent, part, workerID)

		if err != nil {
			log.Printf("Worker %d: failed to write part %d: %v", workerID, part.ID, err)
			return
		}

		fmt.Printf("Worker %d: successfully wrote part %d\n", workerID, part.ID)
	}
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func main() {

	fmt.Println("Starting CSV Cracking process...")

	defer timer("main")()

	csv := service.NewCsv("big_csv.csv")
	numberOfParts := 50

	parts := csv.DivideFileInParts(numberOfParts)

	fmt.Println(parts)

	ch := make(chan service.Part, numberOfParts)

	var wg sync.WaitGroup

	for i := range numberOfParts {
		wg.Add(1)
		go worker(i, ch, &wg, csv)
	}

	for _, part := range parts {
		ch <- part
	}

	close(ch)

	wg.Wait()

}
