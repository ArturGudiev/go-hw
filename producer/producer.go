package producer

import (
	"fmt"
	"os"
	"strings"
)


type Producer struct {
	Filename string
}

func (p Producer) Produce() ([]string, error) {
	data, err := os.ReadFile(p.Filename)
	if err != nil {
		return []string{}, fmt.Errorf("Failed to read file")
	}
	lines := strings.Split(string(data), "\n")
	return lines, nil
}

func NewProducer(filename string) *Producer {
	return &Producer{Filename: filename}
}
