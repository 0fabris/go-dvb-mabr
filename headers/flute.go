package headers

const (
	FLUTE_BASE_HEADER_LEN                    = 16
	FLUTE_TIME_HEADER_SUP_LEN                = 8
	FLUTE_DATA_HEADER_SUP_LEN                = 0
	FLUTE_DATA_XML_DESCRIPTOR_HEADER_SUP_LEN = 20
	FLUTE_CONF_HEADER_SUP_LEN                = 8
	FLUTE_DATA_XML_DESCRIPTOR                = 0x08
)

const (
	FLUTE_DATA_BODY_PACKET byte = 0x03
	FLUTE_DATA_HEAD_PACKET byte = 0x08
)

const (
	NANOCDN_HEADER_BYTE byte = 0x80
)
