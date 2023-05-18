package model

import (
	"context"
	"encoding/json"
	"hdyx/common"
	"hdyx/server"
	"time"
)

type actData struct {
	Uid        int32  `gorm:"uid"`
	ActId      int    `gorm:"actId"`
	TJson      string `gorm:"tJson"`
	CreateTime string `gorm:"createTime"`
	UpdateTime string `gorm:"updateTime"`
}

// GameType : GameName
var actRegister = map[int]func(int32) ActBaseInterface{
	1: NewAct1Model,
}

type RoomIneterface interface {
	gameStart()
}

type ActBaseInterface interface {
	Save() error
	Init() error
}

type playerInterface interface {
	joinRoom()
	leaveRoom()
}

type player struct {
	uid      int32
	playerId int
	power    int
	playTime int64
}

type roomBaseModel struct {
	actType     int
	lvType      int
	roomSize    int
	playerLists []player
	playerNeed  int
	playerNum   int
}

type actBaseModel struct {
	uid     int32
	actId   int
	actInfo map[string]any
}

func (this actBaseModel) Save() error {
	jsonData, _ := json.Marshal(this.actInfo)
	strJsData := string(jsonData)

	now := time.Now()
	strNow := now.Format(time.DateTime)

	actInfo := actData{
		Uid:        this.uid,
		ActId:      this.actId,
		TJson:      strJsData,
		UpdateTime: strNow,
	}

	db := server.GetDb(this.uid, "game")
	db.UpdateData(context.Background(), "act", map[string]any{"uid": this.uid}, actInfo)

	return nil
}

func (this actBaseModel) Init() error {
	db := server.GetDb(this.uid, "game")
	var dest []actData

	err := db.QueryData(context.Background(), "act", map[string]any{"uid": this.uid}, &dest)
	if err != nil {
		common.Logger.ErrorLog(err)
	} else if len(dest) == 0 {
		err = this.actFirstIni()
		return err
	}

	actInfo := dest[0]
	err = json.Unmarshal([]byte(actInfo.TJson), &this.actInfo)

	return err
}

func GetAct(uid int32, actId int) ActBaseInterface {
	fun := actRegister[actId]
	return fun(uid)
}

func (this actBaseModel) actFirstIni() error {
	db := server.GetDb(this.uid, "game")

	jsonData, err := json.Marshal(this.actInfo)
	if err != nil {
		common.Logger.ErrorLog(err)
	}
	strJsData := string(jsonData)

	now := time.Now()
	strNow := now.Format(time.DateTime)

	actInfo := actData{
		Uid:        this.uid,
		ActId:      this.actId,
		TJson:      strJsData,
		CreateTime: strNow,
		UpdateTime: strNow,
	}

	err = db.InsertData(context.Background(), "act", actInfo)

	return err
}
