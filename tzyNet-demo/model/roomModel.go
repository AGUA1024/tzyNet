package model

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	"tzyNet/tCommon"
	"tzyNet/tModel"
	"tzyNet/tNet/ioBuf"
	api "tzyNet/tzyNet-demo/api/protobuf"
	"tzyNet/tzyNet-demo/sdk"
)

const (
	KEY_EXPIRE_TIME = 36000 // reids3600秒过期
)

type roomModel struct {
	ActId          uint32                  `json:"ActId"`
	GameLv         uint32                  `json:"GameLv"`
	IsGame         bool                    `json:"IsGame"`
	PosIdToPlayer  map[uint32]*PlayerModel `json:"PosIdToPlayer"`
	ArrUidAudience []uint64                `json:"ArrUidAudience"`
}

type PlayerModel struct {
	Uid      uint64
	IsMaster bool // 房主
	State    bool // 是否准备
	Head     string
	Name     string
	IsRobot  bool
}

func getRoomModelKey(roomId uint64, actId uint32) string {
	return fmt.Sprintf("room_%d", roomId)
}

func CreateRoom(ctx *tCommon.ConContext, user *sdk.UserInfo, roomId uint64, actId uint32, gameLv uint32) *api.CreateRoom_OutObj {
	redis := tModel.GetCacheById(roomId)

	uid := ctx.GetConGlobalObj().Uid

	roomData := roomModel{
		ActId:  actId,
		GameLv: gameLv,
		PosIdToPlayer: map[uint32]*PlayerModel{
			0: {
				Uid:      uid,
				IsMaster: true,
				State:    false,
				Head:     user.Cover,
				Name:     user.UserName,
			},
		},
		ArrUidAudience: []uint64{},
	}

	data, _ := json.Marshal(roomData)

	ok := redis.RedisWrite(ctx, tModel.REDIS_ROOM, "SETEX", getRoomModelKey(roomId, actId), KEY_EXPIRE_TIME, string(data))
	if !ok {
		return nil
	}

	ret := &api.CreateRoom_OutObj{
		Players: map[uint32]*api.PlayerInfo{
			0: {
				Uid:      uid,
				IsMaster: true,
				State:    false,
				Head:     user.Cover,
				Name:     user.UserName,
			},
		},
	}

	return ret
}

func LoseConnectFunc(ctx *tCommon.ConContext) {
	fmt.Printf("uid:%d,断开连接", ctx.GetConGlobalObj().Uid)
	roomId := ctx.GetConGlobalObj().RoomId
	// 房间不存在
	roomInfo, err := GetGameRoomInfo(ctx, roomId)
	if roomInfo == nil || err != nil {
		fmt.Println("房间不存在")
		return
	}

	// 是否在观众席中
	ok, roomIndex := GetRoomIndex(ctx, roomInfo)
	if ok {
		fmt.Println("离开观众席")
		// 离开观众席位
		roomInfo.ArrUidAudience = append(roomInfo.ArrUidAudience[:roomIndex], roomInfo.ArrUidAudience[roomIndex+1:]...)
	} else {
		fmt.Println("离开游戏坑位")
		playerInfo := roomInfo.PosIdToPlayer
		// 是否在房间
		ok, gameIndex := GetGameIndex(ctx, roomInfo)
		if !ok {
			return
		}

		if isMaster := playerInfo[gameIndex].IsMaster; isMaster {
			// 房主转移
			for i, player := range playerInfo {
				if i == gameIndex || player.IsRobot {
					continue
				}

				// 设置新房主
				if player.IsMaster == false {
					playerInfo[i].IsMaster = true
					break
				}
			}
		}

		// 离开游戏房间
		delete(playerInfo, gameIndex)
		roomInfo.PosIdToPlayer = playerInfo
	}
	fmt.Println("roomInfo.ArrUidAudience:", roomInfo.ArrUidAudience)
	fmt.Println("roomInfo.PosIdToPlayer:", roomInfo.PosIdToPlayer)
	// 离开后如果没有玩家，则销毁房间
	if len(roomInfo.ArrUidAudience) == 0 {
		// 如果房间一个人都没有了，销毁房间、销毁游戏
		if len(roomInfo.PosIdToPlayer) == 0 {
			act1 := GetActModel(ctx, roomInfo.ActId)
			if act1 != nil {
				Destory(ctx, act1)
			}

			DestroyRoom(ctx, roomId)
			return
		} else {
			isDel := true
			for _, player := range roomInfo.PosIdToPlayer {
				if !player.IsRobot {
					isDel = false
					break
				}
			}
			// 如果房间只有机器人,销毁房间、销毁游戏
			if isDel {
				act1 := GetActModel(ctx, roomInfo.ActId)
				if act1 != nil {
					fmt.Println("销毁游戏")
					Destory(ctx, act1)
				}
				fmt.Println("销毁房间")
				DestroyRoom(ctx, roomId)
				return
			}
		}
	}

	// 是否在游戏中
	actModel := GetActModel(ctx, 1)
	if actModel != nil {
		// 获取对局信息
		act1Model := actModel.(*Act1Model)

		// 如果是玩家，则设置掉线托管
		if act1Model.IsPlayer(ctx) {
			act1Model.PlayerLoseConn(ctx)
		}
	}

	// 更新房间数据写入redis
	RoomModelSave(ctx, roomInfo)

	// 广播
	var mpRet = BackGameRoomInfo(roomInfo.PosIdToPlayer)
	broadCastInfo := &api.LeaveGame_OutObj{Players: mpRet}
	MsgRoomBroadcast[*api.LeaveGame_OutObj](ctx, broadCastInfo)
}

func IsRoomExits(ctx *tCommon.ConContext, roomId uint64) (bool, error) {
	redis := tModel.GetCacheById(roomId)

	data, err := redis.RedisQuery("GET", getRoomModelKey(roomId, actId))
	if err != nil {
		return false, err
	}

	return data != nil, nil
}

func GetRoomIndex(ctx *tCommon.ConContext, GameRoomInfo *roomModel) (bool, int) {
	for index, uid := range GameRoomInfo.ArrUidAudience {
		if uid == ctx.GetConGlobalObj().Uid {
			return true, index
		}
	}

	return false, 9999999
}

func GetGameIndex(ctx *tCommon.ConContext, GameRoomInfo *roomModel) (bool, uint32) {
	for index, player := range GameRoomInfo.PosIdToPlayer {
		fmt.Println("player.uid,", player.Uid, "ctx.GetConGlobalObj().uid:", ctx.GetConGlobalObj().Uid)
		if player.Uid == ctx.GetConGlobalObj().Uid {
			return true, index
		}
	}

	return false, 999
}

func GetGameRoomInfo(ctx *tCommon.ConContext, roomId uint64) (*roomModel, error) {
	redis := tModel.GetCacheById(roomId)

	data, err := redis.RedisQuery("GET", getRoomModelKey(roomId, actId))
	if err != nil || data == nil {
		return nil, err
	}

	var roomInfo roomModel
	json.Unmarshal(data.([]byte), &roomInfo)

	return &roomInfo, nil
}

func DestroyRoom(ctx *tCommon.ConContext, roomId uint64) bool {
	redis := tModel.GetCacheById(roomId)

	ok := redis.RedisWrite(ctx, tModel.REDIS_ROOM, "DEL", getRoomModelKey(roomId, actId))
	return ok
}

func MsgRoomBroadcast[T proto.Message](ctx *tCommon.ConContext, obj T) (any, error) {
	// 封装输出数据
	probuf, _ := proto.Marshal(obj)

	out := ioBuf.OutPutBuf{
		Uid:            ctx.GetConGlobalObj().Uid,
		CmdCode:        0,
		ProtocolSwitch: 0,
		CmdMerge:       ctx.GetConGlobalObj().Cmd,
		ResponseStatus: 0,
		Data:           probuf,
	}

	outStream, err := proto.Marshal(&out)
	if err != nil {
		tCommon.Logger.SystemErrorLog("GET_OUT_STREAM_Marshal_ERROR", err)
	}

	// 获取广播列表
	roomInfo, err := GetGameRoomInfo(ctx, ctx.GetConGlobalObj().RoomId)

	if roomInfo == nil || err != nil {
		tCommon.Logger.SystemErrorLog("GetGameRoomInfo_ERR:", err)
	}
	// 玩家广播
	fmt.Println("广播roomInfo.PosIdToPlayer,myuid:", ctx.GetConGlobalObj().Uid)
	for _, player := range roomInfo.PosIdToPlayer {
		fmt.Println("playeruid:", player.Uid)
	}
	fmt.Println("广播roomInfo.ArrUidAudience:", roomInfo.ArrUidAudience)
	for _, player := range roomInfo.PosIdToPlayer {
		uid := player.Uid
		_, ok := tCommon.MpUserStorage[uid]
		if !ok {
			continue
		}

		if err = tCommon.MpUserStorage[uid].WsCon.WriteMessage(tCommon.BinaryMessage, outStream); err != nil {
			continue
		}
	}

	// 观众广播
	for _, uid := range roomInfo.ArrUidAudience {
		_, ok := tCommon.MpUserStorage[uid]
		if !ok {
			continue
		}

		if err = tCommon.MpUserStorage[uid].WsCon.WriteMessage(tCommon.BinaryMessage, outStream); err != nil {
			continue
		}
	}

	return true, nil
}

func RoomModelSave(ctx *tCommon.ConContext, model *roomModel) bool {
	roomId := ctx.GetConGlobalObj().RoomId
	redis := tModel.GetCacheById(roomId)

	data, _ := json.Marshal(model)

	ok := redis.RedisWrite(ctx, tModel.REDIS_ROOM, "SETEX", getRoomModelKey(roomId, actId), KEY_EXPIRE_TIME, string(data))

	return ok
}

func BackGameRoomInfo(roomInfo map[uint32]*PlayerModel) map[uint32]*api.PlayerInfo {
	var mpRet = map[uint32]*api.PlayerInfo{}

	for index, playerInfo := range roomInfo {
		mpRet[index] = &api.PlayerInfo{
			Uid:      playerInfo.Uid,
			State:    playerInfo.State,
			IsMaster: playerInfo.IsMaster,
			Head:     playerInfo.Head,
			Name:     playerInfo.Name,
			IsRobot:  playerInfo.IsRobot,
		}
	}

	return mpRet
}
