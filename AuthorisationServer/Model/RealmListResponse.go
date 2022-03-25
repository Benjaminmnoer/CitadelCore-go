package Model

type RealmListResponse struct {
	Command    uint8
	PacketSize uint16
	Padding    uint32
	RealmCount uint16
	Realms     []RealmInfo
	Wtf        uint8
	Dis        uint8
}

type RealmInfo struct {
	Type       uint8
	Locked     uint8
	Flags      uint8
	Name       string
	Endpoint   string
	Population uint32
	Characters uint8
	Timezone   uint8
	RealmId    uint8
}
