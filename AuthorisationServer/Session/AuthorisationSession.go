package Session

import (
	"CitadelCore/AuthorisationServer/SRP"
	"CitadelCore/Shared/Communication"
)

// Class holding important information about the current auth session.
type AuthorisationSession struct {
	authed         bool
	client         Communication.Client
	AccountName    string
	Srp            *SRP.SRP6
	ReconnectProof [16]byte // Needed for correctly handling reconnect proof. Salt might be more correct name.
}

func StartSession() AuthorisationSession {
	return AuthorisationSession{Srp: SRP.NewSrp()}
}
