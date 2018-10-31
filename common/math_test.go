package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMinInt(t *testing.T) {
  v := MinInt(-1, 1)
  assert.Equal(t, v, -1)
}
