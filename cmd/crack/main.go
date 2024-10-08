package main

import (
	"fmt"
	"os"
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

	csv := service.NewCsv(os.Args[0])

	cli := cli.NewCsvCli(csv)

	cli.Run()

}
