package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/JoseTheodoro42/carck/internal/service"
)

func worker(workerID int, partCh chan service.Part, fileContent []byte, wg *sync.WaitGroup, service *service.Csv) {
	defer wg.Done()
	for part := range partCh {
		err := service.WriteFilePart(fileContent, part, workerID)
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

	csv := service.NewCsv("CSV_NAME.csv")
	numberOfParts := 10

	fileContent, err := csv.ReadFile()

	if err != nil {
		log.Fatalf("Failed to read file: %v\n", err)
	}

	parts := csv.DivideFileInParts(fileContent, numberOfParts)

	ch := make(chan service.Part, numberOfParts)

	var wg sync.WaitGroup

	for i := range numberOfParts {
		wg.Add(1)
		go worker(i, ch, fileContent, &wg, csv)
	}

	for _, part := range parts {
		ch <- part
	}

	close(ch)

	wg.Wait()

}
