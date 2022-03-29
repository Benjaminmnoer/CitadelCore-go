package Model

type LogonChallenge struct {
	Command           uint8
	ProtocolVersion   uint8
	Size              uint16
	Gamename          [4]byte
	Major             uint8
	Minor             uint8
	Patch             uint8
	Build             uint16
	Platform          [4]byte
	Operatingsystem   [4]byte
	Country           [4]byte
	Timezone_bias     uint32
	Ip                uint32
	Accountnamelength uint8
	Accountname       [32]byte // Max 32 character account names. Could be set higher though.
}

func (lc LogonChallenge) GetAccountName() string {
	result := string(lc.Accountname[:lc.Accountnamelength])
	return result
}
