package model

import "encoding/json"

// GameType : GameName
var actRegister = map[int]func(int32) ActBaseInterface{
	1: NewAct1Model,
}

type RoomIneterface interface {
	gameStart()
}

type ActBaseInterface interface {
	save() error
	init() error
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

func (this actBaseModel) save() error {
	_, err := json.Marshal(this.actInfo)

	return err
}

func (this actBaseModel) init() error {
	return nil
}

func GetAct(uid int32, actId int) ActBaseInterface {
	fun := actRegister[actId]
	return fun(uid)
}
