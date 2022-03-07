package Model

type Realm struct {
	Id                   uint64
	Name                 string
	Address              string
	Port                 uint16
	Icon                 uint8
	Flag                 uint8
	Timezone             uint8
	AllowedSecurityLevel uint8
	Population           float32 // Should be unsigned?
	Gamebuild            uint32
}
