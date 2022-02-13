package model

type LoginResult uint8

const (
	Ok LoginResult = iota
	Failed
	Failed2
	Banned
	UnkownAccount
	UnknownAccount3
	AlreadyOnline
	NoTime
	DbBusy
	BadVersion
	DownloadFiled
	Failed3
	Suspended
	Failed4
	Connected
	ParentalControl
	LockedEnforced
)
