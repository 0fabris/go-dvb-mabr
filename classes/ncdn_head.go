package classes

type NCDNPacket interface {
	GetHeader() []byte
	GetData() []byte
}

type NCDNInfoPacket struct {
	NCDNPacket
	RawHeader            []byte  `json:"rawHeader"`
	RawData              []byte  `json:"rawData"`
	Filename             *string `json:"filename"`
	IsRefreshingFilename bool    `json:"isRefreshingFilename"`
}

type NCDNDataPacket struct {
	NCDNPacket
	RawHeader []byte `json:"rawHeader"`
	RawData   []byte `json:"rawData"`
	//Payload           []byte `json:"payload"`
	StartAddress      uint32 `json:"startAddress"`
	IsNewFile         bool   `json:"isNewFile"`
	IsHTTPHeaders     bool   `json:"isHTTPHeaders"`
	IsManifestRelated bool   `json:"isManifestRelated"`
}

type NCDNStream struct {
	URL           string            `json:"url"`
	BroadpeakData string            `json:"bpkData"`
	VideoStreams  map[string]string `json:"videoStreams"`
	VideoPort     uint64            `json:"-"`
	AudioStreams  map[string]string `json:"audioStreams"`
	AudioPort     uint64            `json:"-"`
	DataStreams   map[string]string `json:"dataStreams"`
	DataPort      uint64            `json:"-"`
	ServiceType   string            `json:"serviceType"`
	ServiceID     string            `json:"serviceId"`
	DataSpeed     string            `json:"dataSpeed"`
	RFUn          string            `json:"rfun"`
}
