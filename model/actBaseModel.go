package model

import (
	"context"
	"encoding/json"
	"hdyx/common"
	"hdyx/server"
	"time"
)

type actData struct {
	Uid        uint64 `gorm:"uid"`
	ActId      uint32 `gorm:"actId"`
	TJson      string `gorm:"tJson"`
	CreateTime string `gorm:"createTime"`
	UpdateTime string `gorm:"updateTime"`
}

// GameType : GameName
var actRegister = map[uint32]func(uint64) ActBaseInterface{
	1: NewAct1Model,
}

type ActBaseInterface interface {
	Save() error
	Init() error
}

type actBaseModel struct {
	uid     uint64
	actId   uint32
	actInfo map[string]any
	isOver  bool
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
		common.Logger.SystemErrorLog(err)
	} else if len(dest) == 0 {
		err = this.actFirstIni()
		return err
	}

	actInfo := dest[0]
	err = json.Unmarshal([]byte(actInfo.TJson), &this.actInfo)

	return err
}

func GetAct(uid uint64, actId uint32) ActBaseInterface {
	fun := actRegister[actId]
	return fun(uid)
}

func (this actBaseModel) actFirstIni() error {
	db := server.GetDb(this.uid, "game")

	jsonData, err := json.Marshal(this.actInfo)
	if err != nil {
		common.Logger.SystemErrorLog(err)
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
