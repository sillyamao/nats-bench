package fmt

const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

// Btoa formats bytes with precision set to 2, e.g. "10241B" will be
// "1.01KB". Allowed units are `GB/MB/KB/B` as this function is used
// to format repsonse size etc.
func Btoa(b int64) string {

	basis := Basis{
		Eps: Epsilon{1, []byte("B")},
		Levels: []Level{
			{KB, []byte("KB")},
			{MB, []byte("MB")},
			{GB, []byte("GB")},
		},
	}

	return to2Decimal(b, basis)
}
