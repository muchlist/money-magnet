package slicer_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/muchlist/moneymagnet/foundation/utils/slicer"
)

func TestSlicer(t *testing.T) {
	slice := []int64{24, 26, 27}
	result := slicer.In(26, slice)
	assert.Equal(t, result, true)

	result = slicer.In(25, slice)
	assert.Equal(t, result, false)
}

func TestRemoveFrom(t *testing.T) {
	slice := []int64{24, 26, 27}
	result := slicer.RemoveFrom(26, slice)
	assert.Equal(t, result, []int64{24, 27})
}
