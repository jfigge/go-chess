package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRFtoI(t *testing.T) {
	tests := map[string]struct {
		index uint8
		rank  uint8
		file  uint8
	}{
		"top-left":     {0, 8, 1},
		"top-right":    {7, 8, 8},
		"bottom-left":  {56, 1, 1},
		"bottom-right": {63, 1, 8},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			index := RFtoI(test.rank, test.file)
			assert.Equal(tt, test.index, index)
		})
	}
}

func TestItoRF(t *testing.T) {
	tests := map[string]struct {
		index uint8
		rank  uint8
		file  uint8
	}{
		"top-left":     {0, 8, 1},
		"top-right":    {7, 8, 8},
		"bottom-left":  {56, 1, 1},
		"bottom-right": {63, 1, 8},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			rank, file := ItoRF(test.index)
			assert.Equal(tt, test.rank, rank)
			assert.Equal(tt, test.file, file)
		})
	}
}

func TestRFtoB(t *testing.T) {
	tests := map[string]struct {
		index uint8
		rank  uint8
		file  uint8
	}{
		"top-left":     {63, 8, 1},
		"top-right":    {56, 8, 8},
		"bottom-left":  {7, 1, 1},
		"bottom-right": {0, 1, 8},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			index := RFtoB(test.rank, test.file)
			assert.Equal(tt, test.index, index)
		})
	}
}

func TestBtoRF(t *testing.T) {
	tests := map[string]struct {
		index uint8
		rank  uint8
		file  uint8
	}{
		"top-left":     {63, 8, 1},
		"top-right":    {56, 8, 8},
		"bottom-left":  {7, 1, 1},
		"bottom-right": {0, 1, 8},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			rank, file := BtoRF(test.index)
			assert.Equal(tt, test.rank, rank)
			assert.Equal(tt, test.file, file)
		})
	}
}

func TestNtoRF(t *testing.T) {
	tests := map[string]struct {
		n    string
		rank uint8
		file uint8
		ok   bool
	}{
		"a1": {"a1", 1, 1, true},
		"a8": {"a8", 8, 1, true},
		"h1": {"h1", 1, 8, true},
		"h8": {"h8", 8, 8, true},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			rank, file, valid := NtoRF(test.n)
			assert.Equal(tt, test.rank, rank)
			assert.Equal(tt, test.file, file)
			assert.Equal(tt, test.ok, valid)
		})
	}
}
