package fmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFmt_Btoa(t *testing.T) {

	tests := []struct {
		name string
		b    int64

		// output
		exp string
	}{
		{"Bytes", 100, "100B"},
		{"KB", int64(1.5 * KB), "1.50KB"},
		{"MB", int64(1.5 * MB), "1.50MB"},
		{"GB", int64(1*GB + 512*MB), "1.50GB"},
		{"large-GB", int64(10240 * GB), "10240GB"},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			res := Btoa(test.b)

			assert.Equal(t, test.exp, res)
		})
	}
}
