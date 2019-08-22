package fmt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFmt_Dtoa(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration

		// output
		exp string
	}{

		{"1.50s", 1500 * time.Millisecond, "1.50s"},
		{"60m", 1 * time.Hour, "60m"},
		{"1.50m", 1*time.Minute + 30*time.Second, "1.50m"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := Dtoa(test.d)
			assert.Equal(t, test.exp, res)
		})
	}
}
