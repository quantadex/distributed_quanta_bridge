package coin

import (
	"github.com/magiconair/properties/assert"
	"math"
	"testing"
)

func TestWeiToStellar(t *testing.T) {
	val := WeiToStellar(12344)
	assert.Equal(t, val, int64(0))

	val = WeiToStellar(1000000000000000000)
	assert.Equal(t, val, int64(10000000)) // equal to 1 stellar
}

func TestErc20AmountToStellar(t *testing.T) {
	val := Erc20AmountToStellar(int64(math.Pow10(9)), 9)
	assert.Equal(t, val, int64(10000000))

	val = Erc20AmountToStellar(int64(math.Pow10(18)), 18)
	assert.Equal(t, val, int64(10000000))

	val = Erc20AmountToStellar(1234000, 9)
	assert.Equal(t, val, int64(12340))

	val = Erc20AmountToStellar(10000000, 7)
	assert.Equal(t, val, int64(10000000))
}
