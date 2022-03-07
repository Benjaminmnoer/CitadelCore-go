package Model

type RealmListResponse struct {
	Command    uint8
	PacketSize uint16
	RealmCount uint32
	Realms     []RealmInfo
}

type RealmInfo struct {
	Type       uint8
	Locked     uint8
	Flags      uint8
	Name       string
	Endpoint   string
	Population uint32
	Characters uint8
	RealmId    uint8 // Or uint32. Who knows? Wireshark says 8, Trinity says 32. Lets see what client says
}
