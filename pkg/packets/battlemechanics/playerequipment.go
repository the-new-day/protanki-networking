package battlemechanics

import (
	"github.com/the-new-day/protanki-networking/pkg/codec"
	"github.com/the-new-day/protanki-networking/pkg/codec/complex"
	"github.com/the-new-day/protanki-networking/pkg/packets"
)

// Player equipment in the battle.
// Example json:
//
//	{
//		"battleId": "8337c2cd3515a362",
//		"colormap_id": 3046,
//		"deadColoring": 763,
//		"hull_id": "titan_m3",
//		"turret_id": "smoky_m3",
//		"team_type": "RED",
//		"partsObject": "{some parts objects data here}",
//		"hullResource": 1836,
//		"turretResource": 2190,
//		"sfxData": "<some sfx data here>",
//		"beamTexture": 981,
//		"waveTexture": 983,
//		"sparkTexture": 982,
//		"levelUpSound": 653,
//		"position": {
//			"x": 2379.6272,
//			"y": -3651.7122,
//			"z": 97.50001
//		},
//		"orientation": {
//			"x": 3.8783804E-13,
//			"y": -9.5180135E-15,
//			"z": -0.6118119
//		},
//		"incarnation": 44,
//		"nickname": "ubiza",
//		"state": "active",
//		"maxSpeed": 6,
//		"maxTurnSpeed": 1.5707963267948966,
//		"acceleration": 16,
//		"reverseAcceleration": 17,
//		"sideAcceleration": 20,
//		"turnAcceleration": 1.9198621771937625,
//		"reverseTurnAcceleration": 4.1887902047863905,
//		"mass": 5000,
//		"power": 16,
//		"dampingCoeff": 2100,
//		"turret_turn_speed": 2.1399481958702475,
//		"health": 4158,
//		"rank": 30,
//		"kickback": 2.5,
//		"turretTurnAcceleration": 3.4800119955514934,
//		"impact_force": 3.3,
//		"state_null": false
//	}
type PlayerEquipmentPacket struct {
	packets.BasePacket
}

func NewPlayerEquipmentPacket() *PlayerEquipmentPacket {
	codecs := []codec.Codec{
		codec.Wrap(complex.NewJsonCodec()),
	}

	attributes := []string{
		"json",
	}

	var id int32 = packets.PlayerEquipmentID

	return &PlayerEquipmentPacket{
		BasePacket: *packets.NewBasePacket(id, codecs, attributes),
	}
}

func init() {
	packets.Register(packets.PlayerEquipmentID, "PlayerEquipment", func() packets.Packet {
		return NewPlayerEquipmentPacket()
	})
}
