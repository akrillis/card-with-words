package storage

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestGetAccessToWords(t *testing.T) {
	testData := []struct {
		path  string
		isErr bool
	}{
		{
			path:  "./testDb",
			isErr: false,
		},
		{
			path:  "/0/1/2",
			isErr: true,
		},
	}

	for _, td := range testData {
		_, err := GetAccessToWords(td.path)
		if td.isErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.NoError(t, os.RemoveAll(td.path))
		}
	}
}

func TestDatabase_GetCard(t *testing.T) {
	path := "./testDb"
	words, err := GetAccessToWords(path)
	assert.NoError(t, err)

	for _, i := range []int{2, 4, 8, 16} {
		card, err := words.GetCard(i)
		assert.NoError(t, err)
		assert.Equal(t, i, len(strings.Fields(card)))
	}

	assert.NoError(t, os.RemoveAll(path))
}
