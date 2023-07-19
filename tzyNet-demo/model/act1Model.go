package model

import (
	"fmt"
	"math/rand"
	"time"
	"tzyNet/tCommon"
	"tzyNet/tModel"
	"tzyNet/tNet"
	"tzyNet/tzyNet-demo/app/gateWay/api/protobuf"
)

// 游戏基础配置设置
const (
	actId                            = 1
	maxPlayerNum                     = 5
	minPlayerNum                     = 5
	const_ROBOTCMD_PLAYCARD          = 0x30002
	const_ROBOTCMD_EVENT_HANDLE      = 0x30003
	const_ROBOTCMD_GETCARD_FROM_POOL = 0x30004
	const_ROBOTCMD_TIME_OUT          = 0x30005
)

const (
	const_CARD_SKIP             = iota // 0跳过
	const_CARD_FORBID_1                // 1嫁祸v1
	const_CARD_FORBID_2                // 2嫁祸v2
	const_CARD_REVERSE                 // 3转向
	const_CARD_DRAW_FROM_BOTTOM        // 4抽底
	const_CARD_DEMAND                  // 5索要
	const_CARD_SWAP                    // 6替换
	const_CARD_PREDICT                 // 7预测
	const_CARD_VIEW                    // 8透视
	const_CARD_SHUFFLE                 // 9洗牌
	const_CARD_BOMB                    // 10炸弹
	const_CARD_DISMANTLE               // 11拆除
)

// 牌型对数量的映射
var cfgCardTypeToNum = map[int]int{
	const_CARD_SKIP:             8,
	const_CARD_FORBID_1:         5,
	const_CARD_FORBID_2:         3,
	const_CARD_REVERSE:          5,
	const_CARD_DRAW_FROM_BOTTOM: 3,
	const_CARD_DEMAND:           4,
	const_CARD_SWAP:             3,
	const_CARD_PREDICT:          4,
	const_CARD_VIEW:             4,
	const_CARD_SHUFFLE:          4,
	const_CARD_BOMB:             4,
	const_CARD_DISMANTLE:        6,
}

var (
	EVENT_TYPE_ERROR         uint32 = 99999 // 错误事件
	EVENT_TYPE_SKIP          uint32 = 0     // 跳过
	EVENT_TYPE_FORBID        uint32 = 1     // 嫁祸
	EVENT_TYPE_FORBID_THWICE uint32 = 2     // 嫁祸*2
	EVENT_TYPE_REVERSE       uint32 = 3     // 转向
	EVENT_TYPE_DRAW_BOTTOM   uint32 = 4     // 抽底牌事件
	EVENT_TYPE_STOLEN_TARGET uint32 = 5     // 指定索要目标
	EVENT_TYPE_STOLEN_CARD   uint32 = 6     // 索要事件
	EVENT_TYPE_CARD_SWAP     uint32 = 7     // 交换手牌事件
	EVENT_TYPE_PREDICT       uint32 = 8     // 预测事件
	EVENT_TYPE_CARD_VIEW     uint32 = 9     // 透视三张牌事件
	EVENT_TYPE_SHUFFLE       uint32 = 10    // 洗牌事件
	EVENT_TYPE_BOMB          uint32 = 11    // 炸弹事件
	EVENT_TYPE_DISMANTLE     uint32 = 12    // 拆弹
	EVENT_TYPE_GET_CARD      uint32 = 13    // 摸牌
	EVENT_TYPE_PLAYER_LOSE   uint32 = 14    // 炸弹爆炸，玩家出局
	EVENT_TYPE_GAME_OVER     uint32 = 15    // 游戏结束
	EVENT_TYPE_GAME_INIT     uint32 = 16    // 游戏初始化
	EVENT_TYPE_BOMB_BACK     uint32 = 17    //将炸弹放回牌堆事件
)

type Act1Model struct {
	ActBaseModel
	ActInfo *Act1Info
}

type Act1Info struct {
	PlayerList     map[uint32]*Player
	CurPlayerIndex uint32
	CardPool       []int
	Discarded      []int
	Direction      bool
	BombPlayer     *Player      // 炸弹玩家
	EventPlayer    *EventPlayer // 被索要的玩家
	SeqId          uint32
	MpUidToIndex   map[uint64]uint32
	Rank           []uint64 // 玩家排名,失败后进入队列
	BombNumb       uint32   // 炸弹数量，用于计算炸弹率
}

type EventPlayer struct {
	PlayerIndex uint32
	EventType   uint32
}

type Player struct {
	Uid      uint64 // 玩家id
	Cards    []int  // 手中的牌
	BameNum  uint32 // 是否被嫁祸*2
	IsDie    bool   // 是否已失败
	IsOnline bool   // 是否在线
	IsRobot  bool   // 是否是机器人
}

type RobotEvent struct {
	fun  func([]any)
	args []any
}

func (this RobotEvent) GetFunc() func([]any) {
	return this.fun
}
func (this RobotEvent) GetArgs() []any {
	return this.args
}

func (this *Act1Model) GetActCfg() *GameCfg {
	return &GameCfg{
		ActId:        actId,
		MaxPlayerNum: maxPlayerNum,
		MinPlayerNum: minPlayerNum,
	}
}

// 创建一个新游戏
func (this *Act1Model) NewActModel(ctx *tCommon.ConContext) ActModelInterface {
	// 获取房间信息
	roomInfo, err := GetGameRoomInfo(ctx, ctx.GetConGlobalObj().RoomId)
	if err != nil || roomInfo == nil {
		return nil
	}

	// 抽出炸弹牌，初始牌堆
	deckNoBombs := deckInitNobombs()
	cardPoolShuffle(deckNoBombs)

	// 初始化玩家属性，给玩家发初始手牌
	mpUidToIndex, playerList, cardPool := playerListInit(ctx, deckNoBombs)

	// 发完手牌后加入炸弹牌
	for i := 0; i < cfgCardTypeToNum[const_CARD_BOMB]; i++ {
		cardPool = append(cardPool, const_CARD_BOMB)
	}
	// 洗牌
	cardPoolShuffle(cardPool)

	act1Model := &Act1Model{
		ActBaseModel: ActBaseModel{
			RoomId: ctx.GetConGlobalObj().RoomId,
			ActId:  1,
			IsOver: false,
		},
		ActInfo: &Act1Info{
			PlayerList:     playerList, // playerIndex : player
			CurPlayerIndex: 0,          // 当前出牌玩家索引
			CardPool:       cardPool,   // 卡池
			Discarded:      []int{},    // 弃卡池
			Direction:      true,       // 回合顺序，true 表示正向，false表示反向
			BombPlayer:     nil,        // 手持炸弹的玩家
			EventPlayer:    nil,        // 触发事件的玩家
			SeqId:          0,
			MpUidToIndex:   mpUidToIndex,
			Rank:           []uint64{},                                // 玩家排名
			BombNumb:       uint32(cfgCardTypeToNum[const_CARD_BOMB]), // 炸弹数量
		},
	}

	Save(ctx, act1Model)

	curPlayerIndex := act1Model.ActInfo.CurPlayerIndex
	if act1Model.ActInfo.PlayerList[curPlayerIndex].IsRobot {
		act1Model.RotboToDo(curPlayerIndex)
	}
	return act1Model
}

func (act *Act1Model) GetCardFromPool(ctx *tCommon.ConContext) (uint32, []uint32) {
	// 获取游戏信息
	act1Info := act.ActInfo

	eventType := EVENT_TYPE_GET_CARD
	eventData := []uint32{}

	card := act1Info.CardPool[0]

	act1Info.CardPool = act1Info.CardPool[1:]

	// 如果摸到炸弹
	if card == const_CARD_BOMB {
		eventType = EVENT_TYPE_BOMB
		act1Info.BombPlayer = act1Info.PlayerList[act1Info.CurPlayerIndex]
		// 如果玩家掉线或者是机器人
		act.RotboToDo(act1Info.CurPlayerIndex)
	} else { // 如果不是炸弹则加入手牌，进入下一回合
		// 加入手牌
		curPlayer := act1Info.PlayerList[act1Info.CurPlayerIndex]
		curPlayer.Cards = append(curPlayer.Cards, card)

		// 返回事件内容
		eventData = append(eventData, uint32(card))
		nextTurn(act)
	}

	// 保存数据
	act.ActInfo = act1Info
	Save(ctx, act)

	return eventType, eventData
}

// 对自己使用的卡
func (act *Act1Model) PlayCard(ctx *tCommon.ConContext, cardIndex int) (uint32, []uint32) {
	// 获取游戏信息
	act1Info := act.ActInfo

	// 事件类型与事件内容
	var eventType = EVENT_TYPE_ERROR
	var evnets []uint32

	// 出牌
	card := act.playCard(ctx, cardIndex)
	if card == -1 {
		return EVENT_TYPE_ERROR, nil
	}

	// 不是当前出牌人或者存在事件则报错
	if ctx.GetConGlobalObj().Uid != act1Info.PlayerList[act1Info.CurPlayerIndex].Uid || act1Info.EventPlayer != nil {
		fmt.Println("不是当前出牌人,无法出牌,uid:", ctx.GetConGlobalObj().Uid, "curUid:", act1Info.PlayerList[act1Info.CurPlayerIndex].Uid)
		return EVENT_TYPE_ERROR, nil
	}

	switch card {
	case const_CARD_SKIP: //跳过
		eventType = EVENT_TYPE_SKIP
		// 下一回合
		nextTurn(act)
	case const_CARD_REVERSE: //转向
		eventType = EVENT_TYPE_REVERSE
		act1Info.Direction = !act1Info.Direction
		// 下一回合
		nextTurn(act)
	case const_CARD_DRAW_FROM_BOTTOM: //抽底
		bottonCard := drawBotton(act1Info)
		fmt.Println("抽底牌:", bottonCard)
		if card == const_CARD_BOMB { // 摸到炸弹，不马上进入下一回合
			act1Info.BombPlayer = act1Info.PlayerList[act1Info.CurPlayerIndex]
			eventType = EVENT_TYPE_BOMB
			evnets = append(evnets, uint32(bottonCard))
			// 如果下一位玩家掉线或者是机器人
			act.RotboToDo(act1Info.CurPlayerIndex)
		} else {
			// 将最后一张卡牌加入手牌
			curPlayerIndex := act1Info.CurPlayerIndex
			curPlayerCards := act1Info.PlayerList[curPlayerIndex].Cards
			curPlayerCards = append(curPlayerCards, bottonCard)
			act1Info.PlayerList[curPlayerIndex].Cards = curPlayerCards

			eventType = EVENT_TYPE_DRAW_BOTTOM
			evnets = append(evnets, uint32(bottonCard))
			// 下一回合
			nextTurn(act)
		}
	case const_CARD_PREDICT: //预测
		bombIndex := predictNextBomb(act1Info)
		if bombIndex == -1 {
			fmt.Println("没有找到炸弹")
			return EVENT_TYPE_ERROR, nil
		}
		eventType = EVENT_TYPE_PREDICT
		evnets = append(evnets, uint32(bombIndex))
		// 如果玩家掉线或者是机器人，则继续
		act.RotboToDo(act1Info.CurPlayerIndex)
	case const_CARD_VIEW: //透视
		arrNextThreeCard := getNextThreeCard(act1Info)
		eventType = EVENT_TYPE_CARD_VIEW
		for _, viewCard := range arrNextThreeCard {
			evnets = append(evnets, uint32(viewCard))
		}
		// 如果玩家掉线或者是机器人，则继续
		act.RotboToDo(act1Info.CurPlayerIndex)
		break
	case const_CARD_SHUFFLE: //洗牌
		cardPoolShuffle(act1Info.CardPool)
		eventType = EVENT_TYPE_SHUFFLE
		// 如果玩家掉线或者是机器人，则继续
		act.RotboToDo(act1Info.CurPlayerIndex)
	case const_CARD_DISMANTLE: //拆炸弹
		if act1Info.BombPlayer == nil {
			fmt.Println("没有炸弹要拆")
			return EVENT_TYPE_ERROR, nil
		}
		eventType = EVENT_TYPE_DISMANTLE
		// 触发将炸弹放回牌堆事件
		act1Info.EventPlayer = &EventPlayer{
			PlayerIndex: act1Info.CurPlayerIndex,
			EventType:   EVENT_TYPE_BOMB_BACK,
		}
		act1Info.BombPlayer = nil
		// 如果玩家掉线或者是机器人
		act.RotboToDo(act1Info.CurPlayerIndex)
	case const_CARD_FORBID_1: //嫁祸
		// 返回客户端事件
		eventType = EVENT_TYPE_FORBID

		// 加入选择事件
		act1Info.EventPlayer = &EventPlayer{
			PlayerIndex: act1Info.CurPlayerIndex,
			EventType:   EVENT_TYPE_FORBID,
		}

		// 如果下一位玩家掉线或者是机器人
		act.RotboToDo(act1Info.CurPlayerIndex)
	case const_CARD_FORBID_2: //嫁祸*2
		// 返回客户端事件
		eventType = EVENT_TYPE_FORBID_THWICE
		// 加入选择事件
		act1Info.EventPlayer = &EventPlayer{
			PlayerIndex: act1Info.CurPlayerIndex,
			EventType:   EVENT_TYPE_FORBID_THWICE,
		}

		// 如果下一位玩家掉线或者是机器人
		act.RotboToDo(act1Info.CurPlayerIndex)
	case const_CARD_SWAP: //替换
		// 返回客户端事件
		eventType = EVENT_TYPE_CARD_SWAP
		// 加入选择事件
		act1Info.EventPlayer = &EventPlayer{
			PlayerIndex: act1Info.CurPlayerIndex,
			EventType:   EVENT_TYPE_CARD_SWAP,
		}

		// 如果下一位玩家掉线或者是机器人
		act.RotboToDo(act1Info.CurPlayerIndex)
	case const_CARD_DEMAND: // 打出索要卡
		eventType = EVENT_TYPE_STOLEN_TARGET
		// 被索人触发被索要事件
		act1Info.EventPlayer = &EventPlayer{
			PlayerIndex: act1Info.CurPlayerIndex,
			EventType:   EVENT_TYPE_STOLEN_TARGET,
		}
		// 如果下一位玩家掉线或者是机器人
		act.RotboToDo(act1Info.CurPlayerIndex)
	default:
		fmt.Println("出牌事件错误")
		return EVENT_TYPE_ERROR, nil
	}

	// 保存数据
	act.ActInfo = act1Info
	Save(ctx, act)
	return eventType, evnets
}

// 事件处理
func (act *Act1Model) EventHandler(ctx *tCommon.ConContext, chooseIndex uint32) (uint32, []uint32) {
	// 获取游戏信息
	act1Info := act.ActInfo
	fmt.Println("EventHandler,chooseIndex:", chooseIndex)
	// 事件玩家的信息
	if act1Info.EventPlayer == nil {
		fmt.Println("没有事件可以处理")
		return EVENT_TYPE_ERROR, nil
	}
	eventPlayerIndex := act1Info.EventPlayer.PlayerIndex

	eventType := act1Info.EventPlayer.EventType
	player := act1Info.PlayerList[eventPlayerIndex]

	var events []uint32

	// 被指定的玩家
	if player.IsDie == true {
		fmt.Println("被指定的玩家已经出局了")
		return EVENT_TYPE_ERROR, nil
	}

	switch eventType {
	case EVENT_TYPE_FORBID: // 嫁祸
		// 非法的参数
		player, ok := act1Info.PlayerList[chooseIndex]
		if !ok || player.IsDie {
			fmt.Println("被嫁祸的玩家已经出局了或者被指定的玩家不存在")
			return EVENT_TYPE_ERROR, nil
		}
		// 嫁祸
		eventBlameTarget(act1Info, chooseIndex)

		events = append(events, chooseIndex)
		// 删除事件
		act1Info.EventPlayer = nil
		// 如果下一位玩家掉线或者是机器人
		act.RotboToDo(chooseIndex)
	case EVENT_TYPE_FORBID_THWICE: // 嫁祸*2
		// 非法的参数
		eventPlayer, ok := act1Info.PlayerList[eventPlayerIndex]
		if !ok || eventPlayer.IsDie {
			fmt.Println("被嫁祸*2的玩家已经出局了或者被指定的玩家不存在")
			return EVENT_TYPE_ERROR, nil
		}

		// 嫁祸*2
		eventBlameTargetTwice(act1Info, chooseIndex)
		events = append(events, chooseIndex)

		// 如果被嫁祸的是人机
		act.RotboToDo(act1Info.CurPlayerIndex)

		// 删除事件
		act1Info.EventPlayer = nil
	case EVENT_TYPE_STOLEN_TARGET: // 指定索要对象
		// 非法的参数
		_, ok := act1Info.PlayerList[chooseIndex]
		if !ok || act1Info.PlayerList[eventPlayerIndex].IsDie {
			fmt.Println("索要的玩家不存在")
			return EVENT_TYPE_ERROR, nil
		}
		// 不能对自己使用
		if eventPlayerIndex == chooseIndex {
			fmt.Println("无法对自己索要")
			return EVENT_TYPE_ERROR, nil
		}
		if !eventStolenTarget(act, chooseIndex) {
			fmt.Println("索要的玩家已经出局了")
			return EVENT_TYPE_ERROR, nil
		}
	case EVENT_TYPE_STOLEN_CARD: // 给出索要手牌
		// 给出索要卡
		card, ok := eventStolenCard(act1Info, eventPlayerIndex, int(chooseIndex))
		if !ok {
			return EVENT_TYPE_ERROR, nil
		}

		events = append(events, uint32(card), act1Info.CurPlayerIndex)
		// 下一回合
		nextTurn(act)
		// 删除事件
		act1Info.EventPlayer = nil
	case EVENT_TYPE_CARD_SWAP: // 交换手牌
		// 非法的参数
		if len(act1Info.PlayerList) <= int(chooseIndex) || act1Info.PlayerList[eventPlayerIndex].IsDie {
			return EVENT_TYPE_ERROR, nil
		}
		// 不能对自己使用
		if eventPlayerIndex == chooseIndex {
			return EVENT_TYPE_ERROR, nil
		}
		eventSwapCard(act1Info, chooseIndex)
		events = append(events, chooseIndex)
		// 删除事件
		act1Info.EventPlayer = nil
		// 下一回合
		nextTurn(act)
	case EVENT_TYPE_BOMB_BACK: // 拆弹后将炸弹放回牌堆
		// 删除事件
		act1Info.EventPlayer = nil
		// 下一回合
		nextTurn(act)

		// 传参0~4为将炸弹放回第1张~第5张
		if 0 <= chooseIndex && chooseIndex <= 4 {
			// 炸弹插入牌堆
			backCards := append([]int{}, act1Info.CardPool[chooseIndex:]...)
			act1Info.CardPool = append(act1Info.CardPool[:chooseIndex], const_CARD_BOMB)
			act1Info.CardPool = append(act1Info.CardPool, backCards...)
		} else if chooseIndex == 5 {
			// 炸弹放回牌底
			act1Info.CardPool = append(act1Info.CardPool, const_CARD_BOMB)
		} else if chooseIndex == 6 {
			// 炸弹随机插入牌堆
			var bombIndex int
			if len(act1Info.CardPool) == 0 {
				bombIndex = 0
			} else {
				bombIndex = rand.Intn(len(act1Info.CardPool))
			}

			backCards := append([]int{}, act1Info.CardPool[bombIndex:]...)
			frontCards := append(act1Info.CardPool[:bombIndex], const_CARD_BOMB)
			act1Info.CardPool = append(frontCards, backCards...)
		}
	default:
		return EVENT_TYPE_ERROR, nil
	}

	// 保存数据
	act.ActInfo = act1Info
	Save(ctx, act)

	return eventType, events
}

func deckInitNobombs() []int {
	var deck []int
	for cardType, num := range cfgCardTypeToNum {
		if cardType == const_CARD_BOMB {
			continue
		}

		for i := 0; i < num; i++ {
			deck = append(deck, cardType)
		}
	}

	return deck
}

func playerListInit(ctx *tCommon.ConContext, cardDeck []int) (map[uint64]uint32, map[uint32]*Player, []int) {
	roomInfo, err := GetGameRoomInfo(ctx, ctx.GetConGlobalObj().RoomId)
	if err != nil || roomInfo == nil {
		return nil, nil, nil
	}

	mpRet := map[uint32]*Player{}
	mpUidToIndex := map[uint64]uint32{}

	for index, player := range roomInfo.PosIdToPlayer {
		// 初始化手牌
		var playerCards = cardDeck[:5]
		cardDeck = cardDeck[5:]

		mpUidToIndex[player.Uid] = index

		mpRet[index] = &Player{
			Uid:      player.Uid,
			Cards:    playerCards,
			IsDie:    false,
			BameNum:  0,
			IsOnline: true,
			IsRobot:  player.IsRobot,
		}
	}

	return mpUidToIndex, mpRet, cardDeck
}

// 洗牌
func cardPoolShuffle(cardPool []int) {
	rand.Shuffle(len(cardPool), func(i, j int) {
		cardPool[i], cardPool[j] = cardPool[j], cardPool[i]
	})
}

// 玩家抽牌
func (act *Act1Model) playerDraw(ctx *tCommon.ConContext) int {
	// 获取卡池
	actInfo := act.ActInfo
	cardPool := actInfo.CardPool

	// 抽顶层的卡
	card := cardPool[0]
	cardPool = cardPool[1:]

	// 是否是炸弹,不是炸弹就放入手牌
	if card != const_CARD_BOMB {
		playerList := actInfo.PlayerList
		curIndex := actInfo.CurPlayerIndex

		curPlayerInfo := playerList[curIndex]
		curPlayerInfo.Cards = append(curPlayerInfo.Cards, card)

		actInfo.PlayerList = playerList
	}

	// 更新卡池
	act.ActInfo = actInfo
	Save(ctx, act)

	return card
}

// 玩家打出手中的某张牌
func (act *Act1Model) playCard(ctx *tCommon.ConContext, cardIndex int) int {
	actInfo := act.ActInfo

	// 获取当前出牌玩家
	playerIndex := actInfo.MpUidToIndex[ctx.GetConGlobalObj().Uid]
	player := actInfo.PlayerList[playerIndex]
	// 是否存在这张卡牌
	if cardIndex >= len(player.Cards) || cardIndex < 0 {
		return -1
	}

	// 从玩家手牌中取出要打出的卡牌
	card := player.Cards[cardIndex]

	// 将该卡牌从玩家手牌中删除，并将其加入弃牌堆中
	player.Cards = append(player.Cards[:cardIndex], player.Cards[cardIndex+1:]...)

	// 打出的手牌进入弃卡池
	act.addDiscard(card)
	fmt.Println("出牌:", card)
	return card
}

// 加入弃卡堆
func (act *Act1Model) addDiscard(card int) {
	disCards := act.ActInfo.Discarded
	disCards = append(disCards, card)
}

// 从底抽卡
func drawBotton(actInfo *Act1Info) int {
	// 获取卡池
	cardPool := actInfo.CardPool

	// 获取最后一张牌
	card := cardPool[len(cardPool)-1]

	// 扣除卡池最后一张卡牌
	actInfo.CardPool = cardPool[0 : len(actInfo.CardPool)-1]

	return card
}

// 转到下一个玩家回合
func nextTurn(act *Act1Model) {
	actInfo := act.ActInfo
	playerlist := actInfo.PlayerList
	curIndex := actInfo.CurPlayerIndex

	// 是否被嫁祸*2
	switch playerlist[curIndex].BameNum {
	case 1:
		playerlist[curIndex].BameNum--
	case 2: // 当前还是debuff还剩两回合的时候直接返回,不转到下一个玩家的回合
		// 如果玩家掉线或者是机器人
		act.RotboToDo(curIndex)
		playerlist[curIndex].BameNum--
		return
	}

	fmt.Println("playerList[curIndex].BameNum:", playerlist[curIndex].BameNum)
	var maxPlayerIndex uint32 = 0
	for index, _ := range actInfo.PlayerList {
		if index > maxPlayerIndex {
			maxPlayerIndex = index
		}
	}

	// 转到下一个玩家回合
	if actInfo.Direction == true { // 如果游戏顺序是正向
		curIndex = (curIndex + 1) % (maxPlayerIndex + 1)
	} else { // 如果游戏顺序是反向
		if curIndex == 0 {
			curIndex = maxPlayerIndex
		} else {
			curIndex = (curIndex - 1) % (maxPlayerIndex + 1)
		}
	}

	// 如果下一位玩家死亡，则跳过他的回合
	_, ok := playerlist[curIndex]
	for !ok || playerlist[curIndex].IsDie {
		// 转到下一个玩家回合
		if actInfo.Direction == true { // 如果游戏顺序是正向
			curIndex = (curIndex + 1) % (maxPlayerIndex + 1)
		} else { // 如果游戏顺序是反向
			if curIndex == 0 {
				curIndex = maxPlayerIndex
			} else {
				curIndex = (curIndex - 1) % (maxPlayerIndex + 1)
			}
		}
		_, ok = playerlist[curIndex]
		fmt.Println("curIndex:", curIndex)
	}

	// 如果下一位玩家掉线或者是机器人
	act.RotboToDo(curIndex)

	actInfo.CurPlayerIndex = curIndex
}

// 嫁祸卡
func eventBlameTarget(actInfo *Act1Info, targetIndex uint32) {
	// 被嫁祸的玩家立即到该玩家的回合
	actInfo.CurPlayerIndex = targetIndex
}

// 嫁祸*2卡
func eventBlameTargetTwice(actInfo *Act1Info, targetIndex uint32) {
	// 嫁祸*2 debuff
	actInfo.PlayerList[targetIndex].BameNum = 2
	// 被嫁祸的玩家立即到该玩家的回合
	actInfo.CurPlayerIndex = targetIndex
}

// 玩家索要另一个玩家的手牌
func eventStolenTarget(act *Act1Model, eventPlayerIndex uint32) bool {
	actInfo := act.ActInfo
	eventPlayer := actInfo.PlayerList[eventPlayerIndex]

	if eventPlayer.IsDie {
		return false
	}

	// 玩家有牌才触发给牌事件
	if len(eventPlayer.Cards) > 0 {
		actInfo.EventPlayer = &EventPlayer{
			PlayerIndex: eventPlayerIndex,
			EventType:   EVENT_TYPE_STOLEN_CARD,
		}

		// 如果玩家掉线或者是机器人
		act.RotboToDo(eventPlayerIndex)
	}
	return true
}

// 玩家索要另一个玩家的手牌
func eventStolenCard(actInfo *Act1Info, eventPlayerIndex uint32, cardIndex int) (int, bool) {
	// 被索要玩家
	eventPlayer := actInfo.PlayerList[eventPlayerIndex]
	// 获得卡牌的玩家
	curPlayer := actInfo.PlayerList[actInfo.CurPlayerIndex]
	// 给出的卡
	if len(eventPlayer.Cards) <= cardIndex || eventPlayer.IsDie {
		return 0, false
	}

	card := eventPlayer.Cards[cardIndex]

	// 从对方手牌中移除被索要的卡牌
	eventPlayer.Cards = append(eventPlayer.Cards[:cardIndex], eventPlayer.Cards[cardIndex+1:]...)

	// 增加自己的手牌
	curPlayer.Cards = append(curPlayer.Cards, card) // 将被索要的卡牌加入到自己的手牌中

	return card, true
}

// 预测下一张炸弹
func predictNextBomb(actInfo *Act1Info) int {
	for i, card := range actInfo.CardPool {
		if card == const_CARD_BOMB {
			return i
		}
	}
	return -1 //没有找到炸弹
}

// 玩家使用透视牌，查看接下来三张牌的具体内容
func getNextThreeCard(actInfo *Act1Info) []int {
	cardpool := actInfo.CardPool
	var arrNextThreeCard []int

	num := 0
	for _, card := range cardpool {
		if num == 3 {
			break
		}

		arrNextThreeCard = append(arrNextThreeCard, card)
		num++
	}

	return arrNextThreeCard
}

// swap函数用于两个玩家交换手牌
func eventSwapCard(act1Info *Act1Info, targetPlayIndex uint32) {
	playerList := act1Info.PlayerList
	curIndex := act1Info.CurPlayerIndex

	myInfo := playerList[curIndex]
	targetPlayer := playerList[targetPlayIndex]

	// 交换手牌
	myInfo.Cards, targetPlayer.Cards = targetPlayer.Cards, myInfo.Cards

	// 更新手牌数据
	act1Info.PlayerList = playerList
}

// 拆除炸弹
func (act *Act1Model) disarmBomb(ctx *tCommon.ConContext, bombNewPos int) {
	actInfo := act.ActInfo
	cardpool := actInfo.CardPool

	// 将炸弹放回牌堆
	newCardPool := append(cardpool[:bombNewPos-1], const_CARD_BOMB)
	newCardPool = append(newCardPool, cardpool[bombNewPos+1:]...)
	cardpool = newCardPool

	act.ActInfo = actInfo
	Save(ctx, act)
}

// 超时提醒
func (act *Act1Model) TurnTimeOut(ctx *tCommon.ConContext) (uint32, []uint32) {
	// 获取游戏信息
	act1Info := act.ActInfo

	// 事件类型与事件内容
	var eventType = EVENT_TYPE_ERROR
	var events []uint32

	// 超时序列号+1
	act1Info.SeqId++

	// 是否是拆弹超时
	if act1Info.BombPlayer != nil {
		fmt.Println("被淘汰: ", act1Info.BombPlayer.Uid)
		// 炸弹减少
		act1Info.BombNumb--
		// 玩家出局
		loserIndex := act1Info.MpUidToIndex[act1Info.BombPlayer.Uid]
		act1Info.PlayerList[loserIndex].IsDie = true

		// 出局玩家进入排行榜
		act1Info.Rank = append(act1Info.Rank, act1Info.BombPlayer.Uid)

		// 如果只剩最后一人,则结算
		isGameOver := true
		playerNum := 0
		// 检查是否只剩最后一人
		for _, player := range act1Info.PlayerList {
			if player.IsDie == false {
				playerNum++
			}
			if playerNum >= 2 {
				isGameOver = false
				break
			}
		}

		// 游戏结束结算
		if isGameOver {
			// 幸存者进入排行
			for _, player := range act1Info.PlayerList {
				if player.IsDie == false {
					act1Info.Rank = append(act1Info.Rank, player.Uid)
				}
			}

			// 删除游戏数据
			act.DelAct(ctx)

			for _, playerUid := range act1Info.Rank {
				events = append(events, act1Info.MpUidToIndex[playerUid])
			}
			return EVENT_TYPE_GAME_OVER, events
		}

		// 玩家淘汰
		eventType = EVENT_TYPE_PLAYER_LOSE
		events = append(events, loserIndex)
		act1Info.SeqId++
		// 消除炸弹事件
		act1Info.BombPlayer = nil
		nextTurn(act)
		// 数据落地
		act.ActInfo = act1Info
		Save(ctx, act)
		return eventType, events
	} else if act1Info.EventPlayer != nil { // 超时事件处理
		fmt.Println("EventPlaye.EventType:", act1Info.EventPlayer.EventType)
		switch act1Info.EventPlayer.EventType {
		case EVENT_TYPE_STOLEN_CARD: // 索要选牌事件
			var cardIndex uint32 = 0
			eventType, events = act.EventHandler(ctx, cardIndex)
			fmt.Println("人机处理选牌事件,选择:", cardIndex, "张牌")
			return eventType, events

		case EVENT_TYPE_BOMB_BACK: // 放回炸弹事件
			// 放回随机位置
			eventType, events = act.EventHandler(ctx, 6)
			return eventType, events

		default: // 选人事件
			for key, player := range act1Info.PlayerList {
				if !player.IsDie {
					fmt.Println("人机处理选人事件,选择:", key, "玩家")
					eventType, events = act.EventHandler(ctx, key)
					break
				}
			}
		}

	}

	fmt.Println("超时摸牌")
	// 摸牌
	act1Info.SeqId++
	eventType, events = act.GetCardFromPool(ctx)
	fmt.Println("eventType", eventType, "events", events)
	return eventType, events
}

func RobotHandle(args []any) {
	fmt.Println("人机操作中")
	act := args[0].(*Act1Model)
	uid := args[1].(uint64)
	roomId := args[2].(uint64)
	robotIndex := args[3].(uint32)
	// 创建ctx并初始化
	ctx := tCommon.RegisterConGlobal(9999999)
	if ok := ctx.SetConGlobalUid(uid); !ok {
		return
	}
	if ok := ctx.SetConGlobalRoomId(roomId); !ok {
		return
	}

	// 获取游戏信息
	act1Info := act.ActInfo
	robot := act1Info.PlayerList[robotIndex]

	// 事件类型与事件内容
	var eventType = EVENT_TYPE_ERROR
	var events []uint32

	// 战局分析
	fmt.Println("事件情况：", act1Info.EventPlayer)
	if act1Info.BombPlayer != nil && act1Info.BombPlayer.Uid == robot.Uid {
		fmt.Println("人机处理炸弹....")
		// 炸弹处理
		isBoom := true
		for cardIndex, card := range robot.Cards {
			if card == const_CARD_DISMANTLE {
				fmt.Println("人机拆弹")
				ctx.EventStorageInit(const_ROBOTCMD_PLAYCARD)
				eventType, events = act.PlayCard(ctx, cardIndex)
				isBoom = false
				break
			}
		}
		// 如果不能拆弹
		if isBoom {
			fmt.Println("人机被炸弹淘汰")
			ctx.EventStorageInit(const_ROBOTCMD_TIME_OUT)
			eventType, events = act.TurnTimeOut(ctx)
			if eventType == EVENT_TYPE_GAME_OVER {
				// 数据落地
				tModel.AllRedisSave(ctx)

				// 广播
				outGameInfo := Act1InfoToOutObj(act)

				out := &api.Act1Game_OutObj{
					GameInfo:  outGameInfo,
					EventType: eventType,
					EventData: events,
				}

				MsgRoomBroadcast[*api.Act1Game_OutObj](ctx, out)
				return
			}
		}
	} else if act1Info.EventPlayer != nil && act1Info.EventPlayer.PlayerIndex == robotIndex {
		ctx.EventStorageInit(const_ROBOTCMD_EVENT_HANDLE)
		// 选牌事件
		switch act1Info.EventPlayer.EventType {
		case EVENT_TYPE_STOLEN_CARD:
			var cardIndex uint32 = 0
			eventType, events = act.EventHandler(ctx, cardIndex)
			fmt.Println("人机处理选牌事件,选择:", cardIndex, "张牌")
		case EVENT_TYPE_BOMB_BACK:
			// 放回随机位置
			eventType, events = act.EventHandler(ctx, 6)
		default:
			// 选人事件
			for key, player := range act1Info.PlayerList {
				if !player.IsDie && player.Uid != robot.Uid {
					fmt.Println("人机处理选人事件,选择:", key, "玩家")
					eventType, events = act.EventHandler(ctx, key)
					break
				}
			}
		}
	} else {
		// 自己的回合
		rand.Seed(time.Now().UnixNano())
		robotAct := rand.Intn(2)

		if len(robot.Cards) == 0 {
			// 如果没牌了，则只能摸牌
			robotAct = 0
		} else {
			// 如果没牌可出，则只能摸牌
			canPlayCard := false
			for _, card := range robot.Cards {
				if card != const_CARD_DISMANTLE {
					canPlayCard = true
					break
				}
			}
			if canPlayCard == false {
				robotAct = 0
			}
		}

		switch robotAct {
		case 0: // 摸牌
			ctx.EventStorageInit(const_ROBOTCMD_GETCARD_FROM_POOL)
			fmt.Println("人机摸牌")
			eventType, events = act.GetCardFromPool(ctx)
		case 1: // 出牌
			ctx.EventStorageInit(const_ROBOTCMD_PLAYCARD)
			cardIndex := 0
			for cIndex, card := range robot.Cards {
				if card == const_CARD_DISMANTLE {
					continue
				}
				cardIndex = cIndex
				break
			}
			fmt.Println("人机出牌，第", cardIndex, "张")
			// 随机出牌
			eventType, events = act.PlayCard(ctx, cardIndex)
		}
	}

	// 保存数据
	act.ActInfo = act1Info
	Save(ctx, act)

	// 数据落地
	tModel.AllRedisSave(ctx)

	// 广播
	outGameInfo := Act1InfoToOutObj(act)

	out := &api.Act1Game_OutObj{
		GameInfo:  outGameInfo,
		EventType: eventType,
		EventData: events,
	}

	MsgRoomBroadcast[*api.Act1Game_OutObj](ctx, out)
}

// 数据转化
func Act1InfoToOutObj(act1Model *Act1Model) *api.Act1Info {
	actInfo := act1Model.ActInfo

	// 卡池
	var cardPool []uint32
	for _, card := range actInfo.CardPool {
		cardPool = append(cardPool, uint32(card))
	}

	// 弃卡池
	var discardPool []uint32
	for _, card := range actInfo.Discarded {
		discardPool = append(discardPool, uint32(card))
	}

	var bombPlayer *api.Act1PlayerInfo
	// 玩家信息
	playerList := map[uint32]*api.Act1PlayerInfo{}
	for key, player := range actInfo.PlayerList {
		// 玩家手牌
		var outPlayerCard []uint32
		for _, card := range player.Cards {
			outPlayerCard = append(outPlayerCard, uint32(card))
		}

		// 玩家信息
		playerList[key] = &api.Act1PlayerInfo{
			Uid:     player.Uid,
			Cards:   outPlayerCard,
			IsDie:   player.IsDie,
			BameNum: player.BameNum,
		}

		if actInfo.BombPlayer != nil && actInfo.BombPlayer == player {
			bombPlayer = playerList[key]
		}
	}

	var eventPlayer *api.EventPlayerInfo = nil

	if actInfo.EventPlayer != nil {
		eventPlayer = &api.EventPlayerInfo{
			PlayerIndex: actInfo.EventPlayer.PlayerIndex,
			EventType:   actInfo.EventPlayer.EventType,
		}
	}

	// 排名

	return &api.Act1Info{
		PlayerList:     playerList,
		CurPlayerIndex: actInfo.CurPlayerIndex,
		CardPool:       cardPool,
		DiscardPool:    discardPool,
		Direction:      actInfo.Direction,
		BombPlayer:     bombPlayer,
		EventPlayer:    eventPlayer,
		SeqId:          actInfo.SeqId,
		BombNum:        actInfo.BombNumb,
	}
}

func (act *Act1Model) RotboToDo(playerIndex uint32) {
	actInfo := act.ActInfo
	// 如果下一位玩家掉线或者是机器人
	if actInfo.PlayerList[playerIndex].IsRobot || !actInfo.PlayerList[playerIndex].IsOnline {
		funArgs := []any{act, actInfo.PlayerList[playerIndex].Uid, act.RoomId, playerIndex}
		go func() {
			time.Sleep(3 * time.Second)
			fmt.Println("唤醒机器人")
			tNet.GlobalSysEventChan <- RobotEvent{
				fun:  RobotHandle,
				args: funArgs,
			}
		}()
	}
}

func (act *Act1Model) DelAct(ctx *tCommon.ConContext) {
	modelActId := act.GetActId()
	roomId := act.GetRoomId()

	cache := tModel.GetCacheById(roomId)

	cache.RedisWrite(ctx, tModel.REDIS_ROOM, "DEL", GetActKey(modelActId, roomId))
}

func (act *Act1Model) PlayerLoseConn(ctx *tCommon.ConContext) {
	actInfo := act.ActInfo
	playerIndex := actInfo.MpUidToIndex[ctx.GetConGlobalObj().Uid]
	player := actInfo.PlayerList[playerIndex]
	// 将状态设置为掉线
	player.IsOnline = false
	// 保存数据
	Save(ctx, act)
}

func (act *Act1Model) IsPlayer(ctx *tCommon.ConContext) bool {
	for _, player := range act.ActInfo.PlayerList {
		if player.Uid == ctx.GetConGlobalObj().Uid {
			return true
		}
	}

	return false
}

func (act *Act1Model) PlayerReConn(ctx *tCommon.ConContext) bool {
	uid := ctx.GetConGlobalObj().Uid
	playerIndex, ok := act.ActInfo.MpUidToIndex[uid]
	if !ok {
		return false
	}
	if _, ok = act.ActInfo.PlayerList[playerIndex]; ok {
		return false
	}
	act.ActInfo.PlayerList[playerIndex].IsOnline = true
	Save(ctx, act)
	return true
}
