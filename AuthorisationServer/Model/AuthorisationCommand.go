package model

const (
	AuthLogonChallenge     uint8 = 0x00
	AuthLogonProof               = 0x01
	AuthReconnectChallenge       = 0x02
	AuthReconnectProof           = 0x03
	RealmList                    = 0x10
	TransferInitiate             = 0x30
	TransferData                 = 0x31
	TransferAccept               = 0x32
	TransferResume               = 0x33
	TransferCancel               = 0x34
)
