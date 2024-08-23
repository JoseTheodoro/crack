package service

import (
	"fmt"
	"log"
	"os"
)

type Csv struct {
	fileName string
}

func NewCsv(fileName string) *Csv {
	return &Csv{fileName: fileName}
}

type Part struct {
	ID            int
	startingPoint int
	endingPoint   int
}

func NewPart(ID int, startingPoint int, endingPoint int) *Part {
	return &Part{ID: ID, startingPoint: startingPoint, endingPoint: endingPoint}
}

func (c *Csv) ReadFile() ([]byte, error) {
	fileContent, err := os.ReadFile(c.fileName)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Finished reading file %s\n", c.fileName)
	return fileContent, nil
}

func (c Csv) DivideFileInParts(fileContent []byte, amountOfParts int) []Part {

	fileSize := len(fileContent)
	partSize := fileSize / amountOfParts
	parts := make([]Part, 0, amountOfParts)

	for i := range amountOfParts {
		startingPoint := i * partSize
		endingPoint := startingPoint + partSize

		if i == amountOfParts-1 {
			endingPoint = len(fileContent)
		}

		newPart := NewPart(i, startingPoint, endingPoint)
		parts = append(parts, *newPart)
	}

	fmt.Println("Finished separating part sizes")
	fmt.Printf("File size: %d\n", fileSize)
	fmt.Printf("Part size: %d\n", partSize)

	return parts

}

func (c Csv) WriteFilePart(fileContent []byte, part Part, ID int) error {

	partContent := fileContent[part.startingPoint:part.endingPoint]
	fileName := fmt.Sprintf("csv_part_%d.csv", part.ID)

	fmt.Printf("Worker %d starting to write files!\n", ID)

	err := os.WriteFile(fileName, partContent, 0644)

	if err != nil {
		log.Fatal("Error writing file:", err)
		return err
	}

	fmt.Printf("File written %s\n", fileName)
	return nil
}
