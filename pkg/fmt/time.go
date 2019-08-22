package fmt

import (
	"time"
)

// Dtoa formats duration with precision set to 2, e.g. "10m30.555s" will be
// `10m30.55s". Allowed units are `minute/second/microsecond` as this function
// is used to format request/repsonse time etc. refer Duration.String().
func Dtoa(d time.Duration) string {
	// Largest time is 2540400h10m10.000000000s
	var buf [32]byte
	w := len(buf)

	u := uint64(d)
	neg := d < 0
	if neg {
		u = -u
	}

	var base uint64

	switch {
	case u < uint64(time.Second):
		w--
		buf[w] = 's'
		w--
		buf[w] = 'm'

		base = uint64(time.Millisecond)

	case u < uint64(time.Minute):
		w--
		buf[w] = 's'

		base = uint64(time.Second)

	default:
		// u is greater than time.Minute
		w--
		buf[w] = 'm'

		base = uint64(time.Minute)
	}

	frac := u % base
	u /= base

	w = fmtFrac(buf[:w], frac, base)
	w = fmtInt(buf[:w], u)

	if neg {
		w--
		buf[w] = '-'
	}

	return string(buf[w:])
}

// only keep 2 decimals.
func fmtFrac(buf []byte, v uint64, base uint64) (nw int) {
	w := len(buf)

	u := (v * 100) / base
	if u == 0 {
		return w
	}

	for k := 0; k < 2; k++ {
		digit := u % 10
		w--
		buf[w] = byte(digit) + '0'
		u /= 10
	}

	w--
	buf[w] = '.'

	return w
}

// fmtInt formats v into the tail of buf.
// It returns the index where the output begins.
func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}
