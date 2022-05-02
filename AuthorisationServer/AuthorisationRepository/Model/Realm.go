package Model

type Realm struct {
	Id                   uint8
	Name                 string
	Address              string
	Port                 uint16
	Icon                 uint8
	Flag                 uint8
	Timezone             uint8
	AllowedSecurityLevel uint8
	Population           uint32 // Should be unsigned?
	Gamebuild            uint32
}
