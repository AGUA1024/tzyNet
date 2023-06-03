package sdk

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Sdk struct {
	url     string
	signKey string
}

var GetUserInfoTestSdk = Sdk{
	url:     "http://test5.bugegaming.com/api2/pub/game/getGameUserInfo",
	signKey: "123456",
}

var SendGameRusultTestSdk = Sdk{
	url:     "http://test5.bugegaming.com/api2/pub/game/gameSettlements",
	signKey: "123456",
}

func (this Sdk) GetUrl() string {
	return this.url
}

type GetUserInfo_inBuf struct {
	GameType    uint32   `json:"gameType"`
	Timestamp   int64    `json:"timestamp"`
	Sign        string   `json:"sign"`
	UserNumbers []uint64 `json:"userNumbers"`
}

type GetUserInfoOutBuf struct {
	Code    int                  `json:"code"`
	Data    map[uint64]*UserInfo `json:"data"`
	Msg     string               `json:"msg"`
	TraceId string               `json:"traceId"`
}

type UserInfo struct {
	CodeNum    int    `json:"codeNum"`
	Cover      string `json:"cover"`
	UserId     uint64 `json:"userId"`
	UserName   string `json:"userName"`
	UserNumber string `json:"userNumber"`
}

type SendGameMsgInBuf struct {
	GameType    int            `json:"gameType"`
	GameLv      int            `json:"gameLv"`
	RoomId      int            `json:"roomId"`
	UserRankMap map[uint64]int `json:"userRankMap"`
	Timestamp   int64          `json:"timestamp"`
	Sign        string         `json:"sign"`
}

type SendGameMsgOutBuf struct {
	Code    int    `json:"code"`
	Data    bool   `json:"data"`
	Msg     string `json:"msg"`
	TraceId string `json:"traceId"`
}

// 加密：MD5(gameType + timestamp + "&key=" + signKey)
// 入参：strArgs[0]：gateType，strArgs[1]：timestamp
func (this Sdk) GetSign(strArgs ...string) string {
	gateType := strArgs[0]
	timestamp := strArgs[1]
	sign := this.signKey

	str := gateType + timestamp + "&key=" + sign
	hash := md5.Sum([]byte(str))
	md5Str := hex.EncodeToString(hash[:])
	md5Str = strings.ToUpper(md5Str)

	return md5Str
}

func (this *Sdk) GetPlayerInfoBySdk(actId uint32, uids ...uint64) map[uint64]*UserInfo {
	gameType := actId
	timestamp := time.Now().Unix()
	userNumbers := uids

	sign := this.GetSign(strconv.FormatUint(uint64(gameType), 10), strconv.FormatInt(timestamp, 10))

	req := GetUserInfo_inBuf{
		GameType:    gameType,
		Timestamp:   timestamp,
		Sign:        sign,
		UserNumbers: userNumbers,
	}

	bytesReq, err := json.Marshal(req)
	if err != nil {
		return nil
	}

	byteBack := SdkRequest(this, bytesReq)

	back := GetUserInfoOutBuf{}
	json.Unmarshal(byteBack, &back)

	return back.Data
}

func (this *Sdk) SendGameMsgBySdk(result *SendGameMsgInBuf) *SendGameMsgOutBuf {
	result.Sign = SendGameRusultTestSdk.GetSign(strconv.FormatUint(uint64(result.GameType), 10), strconv.FormatInt(result.Timestamp, 10))

	bytesReq, err := json.Marshal(result)
	if err != nil {
		fmt.Println("error")
	}

	jsonBack := SdkRequest(SendGameRusultTestSdk, bytesReq)

	back := SendGameMsgOutBuf{}
	json.Unmarshal(jsonBack, &back)

	return &back
}
