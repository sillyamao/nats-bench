package fmt

// Level is scale with signs.
type Level struct {
	Threshold uint64

	// If values is greater than current threshold and less than next level's,
	// this sign is used
	Sign []byte
}

type Epsilon struct {
	Base uint64
	Sign []byte
}

// Basis is levels for formatting value of specified type.
// For each type, there is one corresponding basis.
type Basis struct {
	Eps    Epsilon
	Levels []Level
}

// to2Decimal format given input with levels and set precision to 2 decimals.
func to2Decimal(v int64, basis Basis) string {
	// Largest representation is 32 bytes.
	var buf [32]byte
	w := len(buf)

	u := uint64(v)
	neg := v < 0
	if neg {
		u = -u
	}

	base, sign := getBaseAndSign(u, basis)
	frac := u % base
	u /= base

	w = fmtSign(buf[:w], sign)
	w = fmtFrac(buf[:w], frac, base)
	w = fmtInt(buf[:w], u)

	if neg {
		w--
		buf[w] = '-'
	}

	return string(buf[w:])
}

func getBaseAndSign(v uint64, basis Basis) (uint64, []byte) {
	// NOTE: basis must at least have one level.

	lCnt := len(basis.Levels)
	if lCnt == 0 {
		panic("basis must have at least one level")
	}

	fl := basis.Levels[0]
	if v < fl.Threshold {
		return basis.Eps.Base, basis.Eps.Sign
	}

	for k := 0; k < lCnt-1; k++ {
		cur := basis.Levels[k]
		next := basis.Levels[k+1]
		if v >= cur.Threshold &&
			v < next.Threshold {
			return cur.Threshold, cur.Sign
		}
	}

	ll := basis.Levels[lCnt-1]
	return ll.Threshold, ll.Sign
}

func fmtSign(buf []byte, sign []byte) (nw int) {
	w := len(buf)

	for k := len(sign); k > 0; k-- {
		w--
		buf[w] = sign[k-1]
	}

	return w
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
