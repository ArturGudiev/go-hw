package service

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Producer interface {
	Produce() ([]string, error)
}

type Presenter interface {
	Present([]string) error
}

type Service struct {
	prod Producer
	pres Presenter
}

func (s *Service) ReplaceAllLinks(message string) string {
	const httpPrefix = "http://"
	answer := make([]byte, 0, len(message))
	insideLink := false
	index := 0

	for index < len(message) {
		symbol := message[index]
		if insideLink {
			switch symbol {
			case ' ', '\t':
				insideLink = false
				answer = append(answer, symbol)
			default:
				answer = append(answer, '*')
			}
			index++
			continue
		}

		containsPrefix := true
		for i := range httpPrefix {
			if index+i >= len(message) || message[index+i] != httpPrefix[i] {
				containsPrefix = false
				break
			}
		}

		if containsPrefix {
			insideLink = true
			answer = append(answer, httpPrefix...)
			index += len(httpPrefix)
			continue
		}

		answer = append(answer, symbol)
		index++
	}
	return string(answer)
}

func (s *Service) Run() error {
	lines, err := s.prod.Produce()
	if err != nil {
		return err
	}
	modifiedLines := make([]string, len(lines))
	for index, line := range lines {
		modifiedLines[index] = s.ReplaceAllLinks(line)
	}
	err = s.pres.Present(modifiedLines)
	if err != nil {
		return err
	}
	return nil
}

func NewService(prod Producer, pres Presenter) *Service {
	return &Service{
		prod: prod,
		pres: pres,
	}
}

type ProducerImpl struct {
	Filename string
}

func (p ProducerImpl) Produce() ([]string, error) {
	data, err := os.ReadFile(p.Filename)
	if err != nil {
		return []string{}, fmt.Errorf("Failed to read file")
	}
	lines := strings.Split(string(data), "\n")
	return lines, nil
}

func NewProducerImpl(filename string) *ProducerImpl {
	return &ProducerImpl{Filename: filename}
}

type PresenterImpl struct {
	Filename string
}

func (p PresenterImpl) Present(lines []string) error {

	file, err := os.Create("output.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func NewPresenterImpl(filename string) *PresenterImpl {
	return &PresenterImpl{Filename: filename}
}
