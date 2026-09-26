package main

import (
	"flag"
	"fmt"

	"github.com/arturgudiev/go-hw/service"
)


func main() {
	flag.Parse()
	inputFilePath := flag.Arg(0)
	outputFilePath := "output.txt"
	if flag.NArg() >= 2 {
		outputFilePath = flag.Arg(1)
	}

	prod := service.NewProducerImpl(inputFilePath)
	pres := service.NewPresenterImpl(outputFilePath)
	service := service.NewService(prod, pres)
	service.Run()
}
