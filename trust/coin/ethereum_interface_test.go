package coin

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"math"
)

func TestWeiToStellar(t *testing.T) {
	val := WeiToStellar(12344)
	assert.Equal(t, val, int64(0))

	val = WeiToStellar(1000000000000000000)
	assert.Equal(t, int64(val), int64(10000000)) // equal to 1 stellar
}

func TestStellarToWei(t *testing.T) {
	val := StellarToWei(10000000)
	assert.Equal(t, val, uint64(1000000000000000000))
}

func TestWeiToStellarToWei(t *testing.T) {
	m := 1000000000000000000
	x := WeiToStellar(int64(m))
	n := StellarToWei(uint64(x))
	assert.Equal(t, uint64(m), n)
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
