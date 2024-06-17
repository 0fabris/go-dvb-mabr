package classes

func (p *NCDNInfoPacket) GetHeader() []byte {
	return p.RawHeader
}

func (p *NCDNInfoPacket) GetData() []byte {
	return p.RawData
}

func (p *NCDNInfoPacket) GetBlock() []byte {
	return p.RawHeader[14:16]
}

func (p *NCDNDataPacket) GetHeader() []byte {
	return p.RawHeader
}

func (p *NCDNDataPacket) GetData() []byte {
	return p.RawData
}
