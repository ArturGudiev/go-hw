package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	prod := new(MockProducer)
	pres := new(MockPresenter)

	input := []string{"hello http://yandex.com"}
	expected := []string{"hello http://**********"}

	prod.On("Produce").Return(input, nil)
	pres.On("Present", expected).Return(nil)

	svc := NewService(prod, pres)
	svc.Run()

	prod.AssertExpectations(t)
	pres.AssertExpectations(t)

}

func TestService_Produce(t *testing.T) {
	prod := new(MockProducer)
	pres := new(MockPresenter)

	input := []string{"hello http://yandex.com"}
	expected := []string{"hello http://**********"}

	prod.On("Produce").Return(input, nil)
	pres.On("Present", expected).Return(nil)

	svc := NewService(prod, pres)
	svc.Run()

	prod.AssertExpectations(t)
	pres.AssertExpectations(t)
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
	prod := new(MockProducer)
	pres := new(MockPresenter)

	prod.On("Produce").Return(nil, errors.New("Producer error"))

	svc := NewService(prod, pres)
	a.Error(svc.Run())
	pres.AssertNotCalled(t, "Present")
}

func TestService_PresentReturnsError(t *testing.T) {
	a := assert.New(t)
	prod := new(MockProducer)
	pres := new(MockPresenter)

	input := []string{"hello http://yandex.com"}
	expected := []string{"hello http://**********"}

	prod.On("Produce").Return(input, nil)
	pres.On("Present", expected).Return(errors.New("Presenter error"))

	svc := NewService(prod, pres)
	a.Error(svc.Run())
}

func TestProducerImpl_ProducerReturnsLines(t *testing.T) {
	prod := NewProducerImpl("testdata/input.txt")

	lines, err := prod.Produce()

	assert.NoError(t, err)
	assert.Equal(t, []string{"111", "222", "http://google.com"}, lines)
}

func TestProducerImpl_ProduceFileError(t *testing.T) {
	prod := NewProducerImpl("unexisting-file.txt")

	lines, err := prod.Produce()

	assert.Error(t, err)
	assert.Empty(t, lines)
}

func TestPresenterImpl_PresentSucceess(t *testing.T) {
	prod := NewProducerImpl("testdata/input.txt")
	lines, _ := prod.Produce()
	pres := NewPresenterImpl("testdata/output.txt")

	err := pres.Present(lines)
	assert.NoError(t, err)
}
