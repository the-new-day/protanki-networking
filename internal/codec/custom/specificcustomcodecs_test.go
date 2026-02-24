package custom

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper values
func sampleVec(x, y, z float32) map[string]float32 {
	return map[string]float32{"x": x, "y": y, "z": z}
}

func TestBattleLimitsCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewBattleLimitsCodec()

	data := map[string]int32{
		"scoreLimit": int32(1000),
		"timeLimit":  int32(3600),
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data["scoreLimit"], res["scoreLimit"])
	assert.Equal(t, data["timeLimit"], res["timeLimit"])
	assert.Equal(t, 0, buf.Len())
}

func TestRankRangeCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewRankRangeCodec()

	data := map[string]int32{
		"maxRank": int32(15),
		"minRank": int32(5),
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data["maxRank"], res["maxRank"])
	assert.Equal(t, data["minRank"], res["minRank"])
	assert.Equal(t, 0, buf.Len())
}

func TestBattleUserCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewBattleUserCodec()

	data := map[string]any{
		"modLevel": int32(2),
		"deaths":   int16(3),
		"kills":    int16(7),
		"rank":     byte(1),
		"score":    int32(1234),
		"username": "foo",
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data["modLevel"], res["modLevel"])
	assert.Equal(t, data["deaths"], res["deaths"])
	assert.Equal(t, data["kills"], res["kills"])
	assert.Equal(t, data["rank"], res["rank"])
	assert.Equal(t, data["score"], res["score"])
	assert.Equal(t, data["username"], res["username"])
	assert.Equal(t, 0, buf.Len())
}

func TestBattleInfoCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewBattleInfoCodec()

	data := map[string]any{
		"battleID":  "b1",
		"mapName":   "m1",
		"mode":      int32(4),
		"private":   true,
		"proBattle": false,
		"range": map[string]int32{
			"maxRank": int32(20),
			"minRank": int32(10),
		},
		"serverNumber": int32(99),
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data["battleID"], res["battleID"])
	assert.Equal(t, data["mapName"], res["mapName"])
	assert.Equal(t, data["mode"], res["mode"])
	assert.Equal(t, data["private"], res["private"])
	assert.Equal(t, data["proBattle"], res["proBattle"])

	r := res["range"].(map[string]int32)
	assert.Equal(t, int32(20), r["maxRank"])
	assert.Equal(t, int32(10), r["minRank"])

	assert.Equal(t, data["serverNumber"], res["serverNumber"])
	assert.Equal(t, 0, buf.Len())
}

func TestBattleInfoUserCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewBattleInfoUserCodec()

	data := map[string]any{
		"kills":      int32(5),
		"score":      int32(500),
		"suspicious": true,
		"user":       "bar",
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data["kills"], res["kills"])
	assert.Equal(t, data["score"], res["score"])
	assert.Equal(t, data["suspicious"], res["suspicious"])
	assert.Equal(t, data["user"], res["user"])
	assert.Equal(t, 0, buf.Len())
}

func TestBattleNotifierCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewBattleNotifierCodec()

	data := map[string]any{
		"battleInfo": map[string]any{
			"battleID":  "b1",
			"mapName":   "m1",
			"mode":      int32(1),
			"private":   false,
			"proBattle": false,
			"range": map[string]int32{
				"maxRank": int32(1),
				"minRank": int32(0),
			},
			"serverNumber": int32(1),
		},
		"username": "zzz",
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)

	nested := res["battleInfo"].(map[string]any)
	assert.Equal(t, "b1", nested["battleID"])
	assert.Equal(t, "m1", nested["mapName"])
	assert.Equal(t, int32(1), nested["mode"])
	assert.Equal(t, false, nested["private"])
	assert.Equal(t, false, nested["proBattle"])
	assert.Equal(t, int32(1), nested["serverNumber"])
	assert.Equal(t, "zzz", res["username"])
	assert.Equal(t, 0, buf.Len())
}

func TestBattleUserRewardsCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewBattleUserRewardsCodec()

	data := map[string]any{
		"newbiesAbonementBonusReward": int32(10),
		"premiumBonusReward":          int32(20),
		"reward":                      int32(30),
		"userid":                      "u1",
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data, res)
	assert.Equal(t, 0, buf.Len())
}

func TestBattleUserStatsCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewBattleUserStatsCodec()

	data := map[string]any{
		"deaths":   int16(2),
		"kills":    int16(4),
		"score":    int32(200),
		"username": "xyz",
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data, res)
	assert.Equal(t, 0, buf.Len())
}

func TestChatMessageCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewChatMessageCodec()

	author := map[string]any{
		"modLevel": int32(1),
		"ip":       "127.0.0.1",
		"rank":     int32(5),
		"username": "a",
	}
	target := map[string]any{
		"modLevel": int32(0),
		"ip":       "0.0.0.0",
		"rank":     int32(3),
		"username": "b",
	}

	data := map[string]any{
		"authorStatus":  author,
		"systemMessage": false,
		"targetStatus":  target,
		"message":       "hi",
		"warning":       true,
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, author, res["authorStatus"])
	assert.Equal(t, false, res["systemMessage"])
	assert.Equal(t, target, res["targetStatus"])
	assert.Equal(t, "hi", res["message"])
	assert.Equal(t, true, res["warning"])
	assert.Equal(t, 0, buf.Len())
}

func TestFlagInfoCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewFlagInfoCodec()

	data := map[string]any{
		"pole_pos":    sampleVec(1, 2, 3),
		"holder":      "h1",
		"current_pos": sampleVec(4, 5, 6),
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, sampleVec(1, 2, 3), res["pole_pos"])
	assert.Equal(t, "h1", res["holder"])
	assert.Equal(t, sampleVec(4, 5, 6), res["current_pos"])
	assert.Equal(t, 0, buf.Len())
}

func TestMissionRewardCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewMissionRewardCodec()

	data := map[string]any{
		"amount": int32(99),
		"name":   "gold",
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data, res)
	assert.Equal(t, 0, buf.Len())
}

func TestMissionCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewMissionCodec()

	reward1 := map[string]any{
		"amount": int32(1),
		"name":   "a",
	}
	reward2 := map[string]any{
		"amount": int32(2),
		"name":   "b",
	}
	data := map[string]any{
		"freeChange":  true,
		"description": "desc",
		"threshold":   int32(10),
		"image":       int64(123),
		"rewards":     []map[string]any{reward1, reward2},
		"progress":    int32(5),
		"missionID":   int32(7),
		"changeCost":  int32(20),
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data["freeChange"], res["freeChange"])
	assert.Equal(t, data["description"], res["description"])
	assert.Equal(t, data["threshold"], res["threshold"])
	assert.Equal(t, data["image"], res["image"])

	// rewards is a slice of maps
	gotRewards := res["rewards"].([]map[string]any)
	assert.Len(t, gotRewards, 2)
	assert.Equal(t, reward1, gotRewards[0])
	assert.Equal(t, reward2, gotRewards[1])
	assert.Equal(t, data["progress"], res["progress"])
	assert.Equal(t, data["missionID"], res["missionID"])
	assert.Equal(t, data["changeCost"], res["changeCost"])
	assert.Equal(t, 0, buf.Len())
}

func TestMissionStreakCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewMissionStreakCodec()

	data := map[string]any{
		"level":       int32(3),
		"streak":      int32(4),
		"doneToday":   true,
		"questImgID":  int64(10),
		"rewardImgID": int64(20),
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data, res)
	assert.Equal(t, 0, buf.Len())
}

func TestMoveCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewMoveCodec()

	data := map[string]any{
		"angV":        sampleVec(1, 0, 0),
		"control":     byte(5),
		"linV":        sampleVec(0, 1, 0),
		"orientation": sampleVec(0, 0, 1),
		"pos":         sampleVec(9, 8, 7),
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data["angV"], res["angV"])
	assert.Equal(t, data["control"], res["control"])
	assert.Equal(t, data["linV"], res["linV"])
	assert.Equal(t, data["orientation"], res["orientation"])
	assert.Equal(t, data["pos"], res["pos"])
	assert.Equal(t, 0, buf.Len())
}

func TestReferralDataCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewReferralDataCodec()

	data := map[string]any{
		"income":   int32(500),
		"username": "ref",
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data, res)
	assert.Equal(t, 0, buf.Len())
}

func TestTargetHitCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewTargetHitCodec()

	data := map[string]any{
		"direction":     sampleVec(1, 2, 3),
		"localHitPoint": sampleVec(4, 5, 6),
		"numberOfHits":  byte(2),
		"target":        "tgt",
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data, res)
	assert.Equal(t, 0, buf.Len())
}

func TestTargetPositionCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewTargetPositionCodec()

	data := map[string]any{
		"localHitPoint": sampleVec(1, 1, 1),
		"orientation":   sampleVec(2, 2, 2),
		"position":      sampleVec(3, 3, 3),
		"target":        "abc",
		"turretAngle":   float32(45.6),
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data["localHitPoint"], res["localHitPoint"])
	assert.Equal(t, data["orientation"], res["orientation"])
	assert.Equal(t, data["position"], res["position"])
	assert.Equal(t, data["target"], res["target"])
	assert.InDelta(t, data["turretAngle"].(float32), res["turretAngle"].(float32), 0.0001)
	assert.Equal(t, 0, buf.Len())
}

func TestTurretRotateCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewTurretRotateCodec()

	data := map[string]any{
		"angle":   float32(3.14),
		"control": byte(7),
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.InDelta(t, data["angle"].(float32), res["angle"].(float32), 0.0001)
	assert.Equal(t, data["control"], res["control"])
	assert.Equal(t, 0, buf.Len())
}

func TestUserStatusCodec_RoundTrip(t *testing.T) {
	buf := &bytes.Buffer{}
	cc := NewUserStatusCodec()

	data := map[string]any{
		"modLevel": int32(1),
		"ip":       "1.2.3.4",
		"rank":     int32(9),
		"username": "usr",
	}

	n, err := cc.Encode(data, buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)

	res, err := cc.Decode(buf)
	assert.NoError(t, err)
	assert.Equal(t, data, res)
	assert.Equal(t, 0, buf.Len())
}
