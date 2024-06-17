package classes

type MABRFile struct {
	Location    string
	Content     []byte
	HTTPHeaders []byte
}

// Function that converts a []uint8 to uint64
func convertToUInt64(bytes []uint8) uint64 {
	if len(bytes) > 6 {
		return 0
	}

	var result uint64 = 0
	maxPower := (len(bytes) - 1) * 8

	for k, v := range bytes {
		result += uint64(v) << (maxPower - (k * 8))
	}
	return result
}
