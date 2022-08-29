package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunFail(t *testing.T) {
	_, err := NewServer("31212")
	assert.NoError(t, err)
}
