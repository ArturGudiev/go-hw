package service

import "sync"

//go:generate mockery --name=Producer --output=mocks --outpkg=mocks --with-expecter
type Producer interface {
	Produce() ([]string, error)
}

//go:generate mockery --name=Presenter --output=mocks --outpkg=mocks --with-expecter
type Presenter interface {
	Present([]string) error
}

type Service struct {
	prod Producer
	pres Presenter
}

type ReplaceLinksChanData struct {
	index int
	line  string
}

func NewService(prod Producer, pres Presenter) *Service {
	return &Service{
		prod: prod,
		pres: pres,
	}
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

	var wg sync.WaitGroup
	var readerWG sync.WaitGroup
	modifiedLines := make([]string, len(lines))

	tasks := make(chan ReplaceLinksChanData)
	outputChannel := make(chan ReplaceLinksChanData)

	replaceLinksWorker := func(wg *sync.WaitGroup, tasks <-chan ReplaceLinksChanData) {
		for task := range tasks {
			outputChannel <- ReplaceLinksChanData{
				index: task.index,
				line: s.ReplaceAllLinks(task.line),
			}
		}
		wg.Done()
	}

	for range 10 {
		wg.Add(1)
		go replaceLinksWorker(&wg, tasks)
	}

	readerWG.Add(1)
	go func() {
		defer readerWG.Done()
		for taskResult := range outputChannel {
			modifiedLines[taskResult.index] = taskResult.line
		}
	}()

	for index, line := range lines {
		tasks <- ReplaceLinksChanData{index: index, line: line}
	}
	close(tasks)

	wg.Wait()
	close(outputChannel)
	readerWG.Wait()


	err = s.pres.Present(modifiedLines)
	if err != nil {
		return err
	}
	return nil
}
