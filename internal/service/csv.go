package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func (c *Csv) ReadFile(limit Part) ([]byte, error) {

	file, err := os.Open(c.fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bufferSize := limit.endingPoint - limit.startingPoint
	buffer := make([]byte, bufferSize)

	_, err = file.Seek(int64(limit.startingPoint), 0)
	if err != nil {
		return nil, err
	}

	n, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Finished reading %d bytes from file %s\n", n, c.fileName)
	return buffer[:n], nil
}

func (c Csv) DivideFileInParts(amountOfParts int) []Part {

	file, err := os.Open(c.fileName)

	if err != nil {
		fmt.Printf("Could not open file: %v\n", err)
	}

	fileStat, err := file.Stat()

	if err != nil {
		fmt.Printf("Could not get the file stats: %v\n", err)
	}

	fileSize := int(fileStat.Size())
	partSize := fileSize / amountOfParts
	parts := make([]Part, 0, amountOfParts)

	for i := range amountOfParts {
		startingPoint := i * partSize
		endingPoint := startingPoint + partSize

		if i == amountOfParts-1 {
			endingPoint = int(fileStat.Size())
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

	fileNameWithoutPath := filepath.Base(c.fileName)
	fileName := fmt.Sprintf("%s_part_%d.csv", strings.TrimSuffix(fileNameWithoutPath, ".csv"), part.ID)

	fmt.Printf("Worker %d starting to write part %d!\n", ID, part.ID)

	dirName := "Generated Files"
	os.Mkdir(dirName, os.ModePerm)
	err := os.WriteFile(dirName+"/"+fileName, fileContent, 0644)

	if err != nil {
		log.Fatal("Error writing file:", err)
		return err
	}
	fmt.Println("Original file name", c.fileName)
	fmt.Printf("File written %s\n", fileName)
	return nil
}
