package presenter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/arturgudiev/go-hw/producer"
)


func TestPresenter_PresentSucceess(t *testing.T) {
	prod := producer.NewProducer("testdata/input.txt")
	lines, _ := prod.Produce()
	pres := NewPresenter("testdata/output.txt")

	err := pres.Present(lines)
	assert.NoError(t, err)
}
