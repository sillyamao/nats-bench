package fmt

import (
	"time"
)

// Dtoa formats duration with precision set to 2, e.g. "10m30.555s" will be
// `10m30.55s". Allowed units are `minute/second/microsecond` as this function
// is used to format request/repsonse time etc. refer Duration.String().
func Dtoa(d time.Duration) string {
	basis := Basis{
		Eps: Epsilon{uint64(time.Millisecond), []byte("ms")},
		Levels: []Level{
			{uint64(time.Second), []byte("s")},
			{uint64(time.Minute), []byte("m")},
		},
	}

	return to2Decimal(int64(d), basis)
}
