package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/JoseTheodoro42/carck/internal/service"
)

func worker(workerID int, partCh chan service.Part, wg *sync.WaitGroup, service *service.Csv) {
	defer wg.Done()

	for part := range partCh {

		fileContent, err := service.ReadFile(part)
		fmt.Printf("File read successfully by worker: %d\n", workerID)

		service.WriteFilePart(fileContent, part, workerID)

		if err != nil {
			fmt.Printf("Worker %d: failed to write part %d: %v", workerID, part.ID, err)
			return
		}

		fmt.Printf("Worker %d: successfully wrote part %d\n", workerID, part.ID)
	}
}

type CsvCli struct {
	service *service.Csv
}

func NewCsvCli(service *service.Csv) *CsvCli {
	return &CsvCli{service: service}
}

func (cli *CsvCli) Run() {

	fmt.Println("Starting CSV Cracking process...")

	if len(os.Args) < 3 {
		fmt.Println("Usage: crack <file> [parts]")
		return
	}

	commandFile := os.Args[1]
	commandParts, err := strconv.Atoi(os.Args[2])

	if err != nil {
		fmt.Println("Second parameter should be an integer")
	}

	if !strings.HasSuffix(commandFile, ".csv") {
		fmt.Println("File must be of CSV type")
		return
	}

	csv := service.NewCsv(commandFile)
	numberOfParts := commandParts

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
