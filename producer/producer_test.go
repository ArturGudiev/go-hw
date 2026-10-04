package producer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestProducer_ProducerReturnsLines(t *testing.T) {
	prod := NewProducer("testdata/input.txt")

	lines, err := prod.Produce()

	assert.NoError(t, err)
	assert.Equal(t, []string{"111", "222", "http://google.com"}, lines)
}

func TestProducer_ProduceFileError(t *testing.T) {
	prod := NewProducer("unexisting-file.txt")

	lines, err := prod.Produce()

	assert.Error(t, err)
	assert.Empty(t, lines)
}