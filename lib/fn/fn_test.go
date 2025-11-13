package fn_test

import (
	"strconv"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/lib/fn"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	result1 := fn.Map([]int{1, 2, 3, 4}, func(x int) string {
		return "Hello"
	})
	result2 := fn.Map([]int64{1, 2, 3, 4}, func(x int64) string {
		return strconv.FormatInt(x, 10)
	})

	assert.Len(t, result1, 4)
	assert.Len(t, result2, 4)
	assert.Equal(t, []string{"Hello", "Hello", "Hello", "Hello"}, result1)
	assert.Equal(t, []string{"1", "2", "3", "4"}, result2)

	t.Run("given a nil slice, a nil slice of the mapped type is returned", func(t *testing.T) {
		result := fn.Map([]int64(nil), func(x int64) string { return "" })
		assert.Empty(t, result)
		assert.Equal(t, []string(nil), result)
	})

	t.Run("given a zero slice, a zero slice of the mapped type is returned", func(t *testing.T) {
		result := fn.Map([]int64{}, func(x int64) string { return "" })
		assert.Empty(t, result)
		assert.Equal(t, []string{}, result)
	})
}
