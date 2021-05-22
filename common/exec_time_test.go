package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotionTime_AddDate(t *testing.T) {
	a := newNotionTime()
	b := a.AddDate(0, 0, 1)

	assert.False(t, a.time.Equal(b.time))
}
