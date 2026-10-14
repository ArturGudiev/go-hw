package main

import (
	"flag"

	"github.com/arturgudiev/go-hw/service"
	"github.com/arturgudiev/go-hw/producer"
	"github.com/arturgudiev/go-hw/presenter"
)


func main() {
	flag.Parse()
	inputFilePath := flag.Arg(0)
	outputFilePath := "output.txt"
	if flag.NArg() >= 2 {
		outputFilePath = flag.Arg(1)
	}

	prod := producer.NewProducer(inputFilePath)
	pres := presenter.NewPresenter(outputFilePath)
	service := service.NewService(prod, pres)
	service.Run()
}
