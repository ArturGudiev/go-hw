package service

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
	for index, line := range lines {
		modifiedLines[index] = s.ReplaceAllLinks(line)
	}
	err = s.pres.Present(modifiedLines)
	if err != nil {
		return err
	}
	return nil
}

