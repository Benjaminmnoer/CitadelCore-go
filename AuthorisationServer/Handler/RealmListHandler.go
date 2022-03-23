package Handlers

import (
	"CitadelCore/AuthorisationServer/Model"
	"CitadelCore/AuthorisationServer/Repository"
	"fmt"
)

func HandleRealmList(repository Repository.AuthorisationRepository) (Model.RealmListResponse, error) {
	realms, err := repository.GetRealms()

	if err != nil {
		return Model.RealmListResponse{}, fmt.Errorf("Error getting realms\n%s", err)
	}

	result := Model.RealmListResponse{RealmCount: 0, PacketSize: 8, Command: Model.RealmList, Realms: make([]Model.RealmInfo, len(realms)), Wtf: 0x10, Dis: 0}
	for _, realm := range realms {
		var flags uint8 = 0

		realminfo := Model.RealmInfo{
			Type:       realm.Icon,
			Locked:     0,
			Flags:      flags,
			Name:       realm.Name,
			Endpoint:   realm.Address + ":" + fmt.Sprint(realm.Port),
			Population: uint32(realm.Population),
			Characters: 0, //TODO: Add query for this, or add this to already existing query.
			RealmId:    realm.Id,
		}
		result.Realms[result.RealmCount] = realminfo
		result.RealmCount++
		result.PacketSize += uint16(12 + len(realminfo.Name) + len(realminfo.Endpoint))
	}

	return result, nil
}
