package coin

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"math/big"
	"testing"
)

func TestWeiToStellar(t *testing.T) {
	i := new(big.Int)
	fmt.Sscan("12344", i)

	val := WeiToStellar(*i)
	assert.Equal(t, val, int64(0))

	fmt.Sscan("1000000000000000000", i)

	val = WeiToStellar(*i)
	assert.Equal(t, val, int64(10000000)) // equal to 1 stellar
}

func TestErc20AmountToStellar(t *testing.T) {
	i := new(big.Int)
	fmt.Sscan("100000000000", i)
	val := Erc20AmountToGraphene(*i, 9)
	assert.Equal(t, val, int64(10000000))

	fmt.Sscan("100000000000000000000", i)
	val = Erc20AmountToGraphene(*i, 18)
	assert.Equal(t, val, int64(10000000))

	fmt.Sscan("123400000", i)
	val = Erc20AmountToGraphene(*i, 9)
	assert.Equal(t, val, int64(12340))

	fmt.Sscan("1000000000", i)
	val = Erc20AmountToGraphene(*i, 7)
	assert.Equal(t, val, int64(10000000))
}

func TestPowerdelta(t *testing.T) {
	bigInt := big.NewInt(1000000000000000000)
	mul := big.NewInt(10000)
	bigInt.Mul(bigInt, mul)
	result := PowerDelta(*bigInt, 18, 5)
	assert.Equal(t, result, big.NewInt(1000000000))
	result = PowerDelta(*big.NewInt(1000000000000000000), 18, 5)
	assert.Equal(t, result, big.NewInt(100000))

}
