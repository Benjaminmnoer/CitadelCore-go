package Model

type Realm struct {
	Name                 string
	Address              string
	LocalAddress         string
	LocalSubnetMask      string
	Port                 uint16
	Icon                 uint8
	Flag                 uint8
	Timezone             uint8
	AllowedSecurityLevel uint8
	Population           float32 // Should be unsigned?
	Gamebuild            uint32
}
