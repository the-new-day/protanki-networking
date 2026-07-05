package lobby

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Loads detailed battle information.
// Example of a json:
//
//	{
//	  "battleMode": "CTF",
//	  "itemId": "b156433fb81200fb",
//	  "scoreLimit": 10,
//	  "timeLimitInSec": 0,
//	  "preview": 4909,
//	  "maxPeopleCount": 3,
//	  "name": "Sandbox CTF",
//	  "proBattle": true,
//	  "minRank": 30,
//	  "maxRank": 30,
//	  "roundStarted": true,
//	  "spectator": true,
//	  "withoutBonuses": false,
//	  "withoutCrystals": true,
//	  "withoutSupplies": false,
//	  "proBattleEnterPrice": 150,
//	  "timeLeftInSec": -1773519465,
//	  "userPaidNoSuppliesBattle": true,
//	  "proBattleTimeLeftInSec": 1,
//	  "parkourMode": false,
//	  "equipmentConstraintsMode": "NONE",
//	  "reArmorEnabled": true,
//	  "reducedResistance": false,
//	  "esportDropTiming": false,
//	  "withoutGoldBoxes": true,
//	  "withoutGoldSiren": false,
//	  "withoutGoldZone": true,
//	  "withoutMedkit": false,
//	  "withoutMines": false,
//	  "randomGold": true,
//	  "dependentCooldownEnabled": false,
//	  "usersBlue": [
//	    {
//	      "kills": 5,
//	      "score": 176,
//	      "suspicious": false,
//	      "user": "Movement"
//	    },
//	    {
//	      "kills": 5,
//	      "score": 200,
//	      "suspicious": false,
//	      "user": "FinalBoss"
//	    },
//	    {
//	      "kills": 5,
//	      "score": 188,
//	      "suspicious": false,
//	      "user": "Oliver"
//	    }
//	  ],
//	  "usersRed": [
//	    {
//	      "kills": 6,
//	      "score": 150,
//	      "suspicious": false,
//	      "user": "WWF"
//	    },
//	    {
//	      "kills": 7,
//	      "score": 344,
//	      "suspicious": false,
//	      "user": "Jxto"
//	    },
//	    {
//	      "kills": 7,
//	      "score": 157,
//	      "suspicious": false,
//	      "user": "Hot"
//	    }
//	  ],
//	  "scoreRed": 2,
//	  "scoreBlue": 5,
//	  "autoBalance": false,
//	  "friendlyFire": false
//	}
type LoadBattleInfoPacket struct {
	packets.BasePacket
}

func NewLoadBattleInfoPacket() *LoadBattleInfoPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewJsonCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = packets.LoadBattleInfoID

	return &LoadBattleInfoPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.LoadBattleInfoID, "LoadBattleInfo", func() packets.Packet {
		return NewLoadBattleInfoPacket()
	})
}
