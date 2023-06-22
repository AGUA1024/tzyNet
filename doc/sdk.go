package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SdkInterface interface {
	GetSign(strArgs ...string) string
	Request(jsonReqBuf string) string
}

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

type reqBuf struct {
	GameType    uint32   `json:"gameType"`
	Timestamp   int64    `json:"timestamp"`
	Sign        string   `json:"sign"`
	UserNumbers []uint64 `json:"userNumbers"`
}

type getUserInfoBackBuf struct {
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

type SendGameMsgRet struct {
	Code    int    `json:"code"`
	Data    bool   `json:"data"`
	Msg     string `json:"msg"`
	TraceId string `json:"traceId"`
}

// 加密：MD5(gameType + timestamp + "&key=" + signKey)
// 入参：strArgs[0]：gateType，strArgs[1]：timestamp
func (this *Sdk) GetSign(strArgs ...string) string {
	gateType := strArgs[0]
	timestamp := strArgs[1]
	sign := this.signKey

	str := gateType + timestamp + "&key=" + sign
	hash := md5.Sum([]byte(str))
	md5Str := hex.EncodeToString(hash[:])
	md5Str = strings.ToUpper(md5Str)

	return md5Str
}

type gameResult struct {
	GameType    int            `json:"gameType"`
	GameLv      int            `json:"gameLv"`
	RoomId      int            `json:"roomId"`
	UserRankMap map[uint64]int `json:"userRankMap"`
	Timestamp   int64          `json:"timestamp"`
	Sign        string         `json:"sign"`
}

func (this *Sdk) Request(jsonReqBuf string) []byte {
	byteData := []byte(jsonReqBuf)
	fmt.Println(this.url)
	req, err := http.NewRequest(http.MethodPost, this.url, bytes.NewReader(byteData))
	if err != nil {
		//common.Logger.GameErrorLog(ctx, common.ERR_PLAYERINFO_REQUEST, "建立http会话失败")
		return nil
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//common.Logger.GameErrorLog(ctx, common.ERR_PLAYERINFO_REQUEST, "向不鸽平台请求玩家信息失败")
		return nil
	}
	defer resp.Body.Close()

	// 解析响应
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Request error")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	return body
}

func (this *Sdk) GetPlayerInfoBySdk(actId uint32, uids ...uint64) *map[uint64]*UserInfo {
	gameType := actId
	timestamp := time.Now().Unix()
	userNumbers := uids

	sign := GetUserInfoTestSdk.GetSign(strconv.FormatUint(uint64(gameType), 10), strconv.FormatInt(timestamp, 10))

	req := reqBuf{
		GameType:    gameType,
		Timestamp:   timestamp,
		Sign:        sign,
		UserNumbers: userNumbers,
	}

	bytesReq, err := json.Marshal(req)
	if err != nil {
		fmt.Println("error")
	}

	jsonReq := string(bytesReq)

	jsonBack := GetUserInfoTestSdk.Request(jsonReq)

	back := getUserInfoBackBuf{}
	json.Unmarshal(jsonBack, &back)

	return &back.Data
}

func (this *Sdk) SendGameMsgBySdk(result *gameResult) *SendGameMsgRet {
	result.Sign = SendGameRusultTestSdk.GetSign(strconv.FormatUint(uint64(result.GameType), 10), strconv.FormatInt(result.Timestamp, 10))

	bytesReq, err := json.Marshal(result)
	if err != nil {
		fmt.Println("error")
	}

	jsonReq := string(bytesReq)

	jsonBack := SendGameRusultTestSdk.Request(jsonReq)

	back := SendGameMsgRet{}
	json.Unmarshal(jsonBack, &back)

	return &back
}

func main() {
	sgin := SendGameRusultTestSdk.GetSign(strconv.FormatUint(1, 10), strconv.FormatUint(1685611407000, 10))

	gr := &gameResult{
		GameType: 1,
		GameLv:   5,
		RoomId:   444481,
		UserRankMap: map[uint64]int{
			1160707389: 2,
		},
		Timestamp: 1685611407000,
		Sign:      sgin,
	}

	back := SendGameRusultTestSdk.SendGameMsgBySdk(gr)
	fmt.Println(back)
}
