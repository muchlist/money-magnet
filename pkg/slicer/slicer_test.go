package slicer_test

import (
	"testing"

	"github.com/muchlist/moneymagnet/pkg/slicer"
	"github.com/stretchr/testify/assert"
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

func TestCsvToSliceInt(t *testing.T) {
	tests := map[string]struct {
		input      string
		wantResult []int
		wantError  bool
	}{
		"all input int":       {input: "1,2,3,4,5,1,2", wantResult: []int{1, 2, 3, 4, 5, 1, 2}, wantError: false},
		"any non convertable": {input: "1,2,muchlis,3", wantResult: []int{1, 2, 3}, wantError: true},
		"all non convertabe":  {input: "muchls", wantResult: []int{}, wantError: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := slicer.CsvToSliceInt(tc.input)
			assert.Equal(t, tc.wantResult, result, "wrong result")
			assert.Equal(t, tc.wantError, err != nil)
		})
	}
}
