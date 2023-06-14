package common

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"runtime"
	"strconv"
)

// connectGorutine局部空间
var MpConRoutineStorage = map[uint64]*ConGlobalStorage{}

var MpUserStorage = map[uint64]*userModel{}

type userModel struct {
	WsCon  *websocket.Conn
	RoomId uint64
}

type ConGlobalStorage struct {
	WsCon        *websocket.Conn
	Uid          uint64
	Cmd          uint32
	RoomId       uint64
	EventStorage *eventGlobalStorage
}

type eventGlobalStorage struct {
	playerCache []CacheEvent
	roomCache   []CacheEvent
}

type CacheEvent interface {
	GetCommand() string
	GetArgs() []interface{}
}

type CacheOperator interface {
	CacheSave(arrCacheEvent []CacheEvent) error
}

// connectGorutine上下文
type ConContext struct {
	ConnectId uint64
}

// 事件全局存储空间初始化
func (ctx *ConContext) EventStorageInit(cmd uint32) {
	MpConRoutineStorage[ctx.GetConnectId()].Cmd = cmd
	MpConRoutineStorage[ctx.GetConnectId()].EventStorage = &eventGlobalStorage{
		playerCache: []CacheEvent{},
		roomCache:   []CacheEvent{},
	}
}

func (ctx *ConContext) PlayerRedisEventPush(playerRedisEvent CacheEvent) bool {
	global := ctx.GetConGlobalObj()
	if global == nil || global.EventStorage.playerCache == nil || global.EventStorage.playerCache == nil {
		return false
	}
	global.EventStorage.playerCache = append(global.EventStorage.playerCache, playerRedisEvent)

	return true
}

func (ctx *ConContext) RoomRedisEventPush(roomRedisEvent CacheEvent) bool {
	global := ctx.GetConGlobalObj()
	if global == nil || global.EventStorage == nil || global.EventStorage.playerCache == nil || global.EventStorage.playerCache == nil {
		return false
	}
	global.EventStorage.roomCache = append(global.EventStorage.roomCache, roomRedisEvent)

	return true
}

// 所有内存数据落地缓存
func (ctx *ConContext) AllCacheSave(playerCacheOp, RoomCacheOp CacheOperator) error {

	// 获取redis事件
	roomCache := ctx.GetConGlobalObj().EventStorage.roomCache
	playerCache := ctx.GetConGlobalObj().EventStorage.playerCache

	fmt.Println("roomCache:")
	fmt.Println(roomCache)
	// 本地内存落地redis
	err := playerCacheOp.CacheSave(roomCache)
	err = RoomCacheOp.CacheSave(playerCache)

	return err
}

func (this *ConContext) GetConnectId() uint64 {
	return this.ConnectId
}

// readOnly
func (ctx *ConContext) GetConGlobalObj() *ConGlobalStorage {
	conGlobalObj, ok := MpConRoutineStorage[ctx.GetConnectId()]
	if !ok {
		return nil
	}
	return conGlobalObj
}

func (ctx *ConContext) SetConGlobalWsCon(conn *websocket.Conn) bool {
	conId := ctx.GetConnectId()

	conGlobalObj, ok := MpConRoutineStorage[conId]
	if !ok {
		return ok
	}
	conGlobalObj.WsCon = conn

	MpConRoutineStorage[conId] = conGlobalObj

	return true
}

func (ctx *ConContext) SetConGlobalCmd(cmd uint32) bool {
	conId := ctx.GetConnectId()

	conGlobalObj, ok := MpConRoutineStorage[conId]
	if !ok {
		return ok
	}
	conGlobalObj.Cmd = cmd

	MpConRoutineStorage[conId] = conGlobalObj

	return true
}

func (ctx *ConContext) SetConGlobalRoomId(roomId uint64) bool {
	conId := ctx.GetConnectId()

	conGlobalObj, ok := MpConRoutineStorage[conId]
	if !ok {
		return ok
	}

	conGlobalObj.RoomId = roomId

	MpConRoutineStorage[conId] = conGlobalObj

	return true
}

func (ctx *ConContext) SetConGlobalUid(uid uint64) bool {
	conId := ctx.GetConnectId()

	conGlobalObj, ok := MpConRoutineStorage[conId]
	if !ok {
		return ok
	}
	conGlobalObj.Uid = uid

	MpConRoutineStorage[conId] = conGlobalObj

	return true
}

func (ctx *ConContext) RegisterUserForStorage() bool {
	if MpUserStorage == nil || ctx.GetConGlobalObj() == nil {
		return false
	}
	uid := ctx.GetConGlobalObj().Uid
	wsCon := ctx.GetConGlobalObj().WsCon
	roomId := ctx.GetConGlobalObj().RoomId

	MpUserStorage[uid] = &userModel{
		WsCon:  wsCon,
		RoomId: roomId,
	}

	return true
}

func newConContext() *ConContext {
	return &ConContext{
		ConnectId: getGoroutineID(),
	}
}

// 注册connectGorutine全局空间
func RegisterConGlobal() *ConContext {
	ctx := newConContext()
	conId := ctx.GetConnectId()
	//if _, ok := MpConRoutineStorage[conId]; ok {
	//Logger.SystemErrorLog("RoutineStorage_MEM_OVERWRITE")
	//}

	MpConRoutineStorage[conId] = &ConGlobalStorage{
		WsCon:        nil,
		Uid:          0,
		Cmd:          0,
		RoomId:       0,
		EventStorage: nil,
	}

	return ctx
}

func getGoroutineID() uint64 {
	arrByte := make([]byte, 64)
	arrByte = arrByte[:runtime.Stack(arrByte, false)]
	arrByte = bytes.TrimPrefix(arrByte, []byte("goroutine "))
	arrByte = arrByte[:bytes.IndexByte(arrByte, ' ')]
	GoroutineID, _ := strconv.ParseUint(string(arrByte), 10, 64)
	return GoroutineID
}
