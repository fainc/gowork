package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAmountTrans(t *testing.T) {
	t.Run("Int642Yuan", func(t *testing.T) {
		f, s := AmountTrans().Int642Yuan(1, 3, 2)
		fmt.Println("f：", f)
		fmt.Println("s：", s)
		assert.Equal(t, f, 0.33)
	})
	t.Run("Yuan2Int64", func(t *testing.T) {
		i := AmountTrans().Yuan2Int64(2.3, 100)
		fmt.Println("i：", i)
		assert.Equal(t, i, int64(230))
	})
	t.Run("Yuan2Int", func(t *testing.T) {
		i := AmountTrans().Yuan2Int(2.3, 100)
		fmt.Println("i：", i)
		assert.Equal(t, i, 230)
	})
}
