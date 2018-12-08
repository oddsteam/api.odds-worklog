package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCSVFile(t *testing.T) {
	file, filename, err := CreateCVSFile("corporate")
	assert.NotNil(t, file)
	assert.NotEmpty(t, filename)
	assert.NoError(t, err)

	// remove file after test
	os.Remove(filename)
}
