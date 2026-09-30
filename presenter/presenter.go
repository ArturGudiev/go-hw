package presenter

import (
	"bufio"
	"os"
)


type Presenter struct {
	Filename string
}


func NewPresenter(filename string) *Presenter {
	return &Presenter{Filename: filename}
}

func (p Presenter) Present(lines []string) error {

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
