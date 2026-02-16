package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelVendorCode(t *testing.T) {
	t.Run("test large index vendor code", func(t *testing.T) {
		vc := VendorCode{index: 381}
		assert.Equal(t, "AOR", vc.String())
	})
}
