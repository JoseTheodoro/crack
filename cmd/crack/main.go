package main

import (
	"fmt"
	"time"

	"github.com/JoseTheodoro42/carck/internal/cli"
	"github.com/JoseTheodoro42/carck/internal/service"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func main() {

	defer timer("main")()

	csv := service.NewCsv("big_csv.csv")

	cli := cli.NewCsvCli(csv)

	cli.Run()

}
