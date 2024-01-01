package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignVerify(t *testing.T) {
	kg := NewKeyGen()
	payload := []byte("hello world")

	assert.NoError(t, kg.Sign(payload))
	assert.True(t, kg.Verify(payload))
}
