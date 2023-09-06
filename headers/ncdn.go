package headers

// mABR Settings
const (
	NCDN_PACKET_SIZE      = 2048
	NCDN_HEADER_LEN_HPACK = 16
	NCDN_HEADER_LEN_DPACK = 20
	NCDN_FILENAME_MAX_LEN = 300
	NCDN_INFO_DATA_SEP    = ";"
)

// Consts
var NCDN_HEADER_HEAD_PACKET []byte = []byte{0x80, 0xA1}

var NCDN_HEADER_DATA_PACKET []byte = []byte{0x80, 0x21}
