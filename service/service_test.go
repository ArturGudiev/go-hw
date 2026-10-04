package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/arturgudiev/go-hw/service/mocks"
)

type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Produce() ([]string, error) {
	args := m.Called()
	lines, _ := args.Get(0).([]string)
	return lines, args.Error(1)
}

type MockPresenter struct {
	mock.Mock
}

func (m *MockPresenter) Present(lines []string) error {
	args := m.Called(lines)
	return args.Error(0)
}

func TestService_Run(t *testing.T) {
	prod := mocks.NewProducer(t)
	pres := mocks.NewPresenter(t)

	input := []string{"hello http://yandex.com"}
	expected := []string{"hello http://**********"}

	prod.On("Produce").Return(input, nil)
	pres.On("Present", expected).Return(nil)

	svc := NewService(prod, pres)
	svc.Run()

}

func TestService_Produce(t *testing.T) {
	prod := mocks.NewProducer(t)
	pres := mocks.NewPresenter(t)


	input := []string{"hello http://yandex.com"}
	expected := []string{"hello http://**********"}

	prod.On("Produce").Return(input, nil)
	pres.On("Present", expected).Return(nil)

	svc := NewService(prod, pres)
	svc.Run()

}

func TestService_ReplaceAllLinks(t *testing.T) {
	a := assert.New(t)
	svc := NewService(nil, nil)
	checks := [][2]string{
		{"", ""},
		{"http://yandex.com", "http://**********"},
		{"https://google.com", "https://google.com"},
		{"just text", "just text"},
		{"http://x.com\tnext", "http://*****\tnext"},
	}
	for _, check := range checks {

		a.Equal(check[1], svc.ReplaceAllLinks(check[0]))

	}
}

func TestService_ProduceReturnsError(t *testing.T) {
	a := assert.New(t)
	prod := mocks.NewProducer(t)
	pres := mocks.NewPresenter(t)

	prod.On("Produce").Return(nil, errors.New("Producer error"))

	svc := NewService(prod, pres)
	a.Error(svc.Run())
	pres.AssertNotCalled(t, "Present")
}

func TestService_PresentReturnsError(t *testing.T) {
	a := assert.New(t)
	prod := mocks.NewProducer(t)
	pres := mocks.NewPresenter(t)

	input := []string{"hello http://yandex.com"}
	expected := []string{"hello http://**********"}

	prod.On("Produce").Return(input, nil)
	pres.On("Present", expected).Return(errors.New("Presenter error"))

	svc := NewService(prod, pres)
	a.Error(svc.Run())
}


