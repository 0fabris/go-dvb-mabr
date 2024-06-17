package flute

// build a (LSB) mask from a given number of bits
func build32bitLSBMask(bits int) uint32 {
	var m uint32
	for i := 0; i < bits; i++ {
		m |= (1 << i)
	}
	return m
}
