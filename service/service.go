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

	modifiedLines := make([]string, len(lines))
	maxWorkers := min(len(lines), 10)

	tasksChan := make(chan ReplaceLinksChanData, maxWorkers)
	resultsChan := make(chan ReplaceLinksChanData, maxWorkers)
	sem := make(chan struct{}, maxWorkers)
	wg := &sync.WaitGroup{}
	readerWg := &sync.WaitGroup{}

	go func() {
		for i, line := range lines {
			tasksChan <- ReplaceLinksChanData{index: i, line: line}
		}
		close(tasksChan)
	}()

	go func () {
		for task := range tasksChan {
			wg.Add(1)
			sem <- struct{}{}
			go func() {
				defer wg.Done()
				defer func() {<- sem}()
				resultsChan <- ReplaceLinksChanData{
					index: task.index,
					line: s.ReplaceAllLinks(task.line),
				}	
			}()
		}
		wg.Wait()
		close(resultsChan)
	}()

	readerWg.Add(1)
	go func() {
		defer readerWg.Done()
		for taskResult := range resultsChan {
			modifiedLines[taskResult.index] = taskResult.line
		}	
	}()
	readerWg.Wait()

	err = s.pres.Present(modifiedLines)
	if err != nil {
		return err
	}
	return nil
}
