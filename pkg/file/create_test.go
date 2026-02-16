package file

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFile(t *testing.T) {
	file, filename, err := CreateFile("corporate")
	assert.NotNil(t, file)
	assert.NotEmpty(t, filename)
	assert.NoError(t, err)

	// remove file after test
	os.Remove(filename)
}
