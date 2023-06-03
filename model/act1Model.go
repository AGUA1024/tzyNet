package model

import (
	"hdyx/common"
	"math/rand"
)

// 游戏基础配置设置
const (
	actId        = 1
	maxPlayerNum = 5
	minPlayerNum = 5
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

type Act1Model struct {
	actBaseModel
}

type Player struct {
	Uid    uint64 //玩家id
	Cards  []int  // 手中的牌
	Status bool   // 是否已失败
}

func (this *Act1Model) GetActCfg() *GameCfg {
	return &GameCfg{
		ActId:        actId,
		MaxPlayerNum: maxPlayerNum,
		MinPlayerNum: minPlayerNum,
	}
}

// 创建一个新游戏
func (this Act1Model) NewActModel(ctx *common.ConContext) ActBaseInterface {
	// 获取房间信息
	roomInfo, err := GetGameRoomInfo(ctx, ctx.GetConGlobalObj().RoomId)
	if err != nil || roomInfo == nil {
		return nil
	}

	// 抽出炸弹牌，初始牌堆
	deckNoBombs := deckInitNobombs()
	cardPoolShuffle(deckNoBombs)

	// 初始化玩家属性，给玩家发初始手牌
	playerList, cardPool := playerListInit(ctx, deckNoBombs)

	// 发完手牌后加入炸弹牌
	for i := 0; i < cfgCardTypeToNum[const_CARD_BOMB]; i++ {
		cardPool = append(cardPool, const_CARD_BOMB)
	}
	// 洗牌
	cardPoolShuffle(cardPool)

	act1Model := &Act1Model{
		actBaseModel: actBaseModel{
			RoomId: ctx.GetConGlobalObj().RoomId,
			ActId:  1,
			IsOver: false,
			ActInfo: map[string]any{
				"playerList":       playerList, // playerId : cardId
				"curPlayerId":      0,          // 当前出牌玩家id
				"curPlayerRemTime": -1,         // 当前玩家剩余时间
				"cardPool":         cardPool,   // 卡池
				"discarded":        []int{},    // 弃卡池
				"direction":        true,       // 回合顺序，true 表示正向，false表示反向
			},
		},
	}

	this.Save(ctx)
	return act1Model
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

func playerListInit(ctx *common.ConContext, cardDeck []int) (*map[uint32]Player, []int) {
	roomInfo, err := GetGameRoomInfo(ctx, ctx.GetConGlobalObj().RoomId)
	if err != nil || roomInfo == nil {
		return nil, nil
	}

	mpRet := map[uint32]Player{}
	for index, player := range roomInfo.PosIdToPlayer {
		// 初始化手牌
		var playerCards = cardDeck[:5]
		cardDeck = cardDeck[5:]

		mpRet[index] = Player{
			Uid:    player.Uid,
			Cards:  playerCards,
			Status: false,
		}
		index++
	}

	return &mpRet, cardDeck
}

// 洗牌
func cardPoolShuffle(cardpool []int) {
	rand.Shuffle(len(cardpool), func(i, j int) {
		cardpool[i], cardpool[j] = cardpool[j], cardpool[i]
	})
}

//// // 根据玩家数量初始化游戏
//func NewGame(playerCount int, playerId []int16) *Game {
//	// 创建玩家数组，并初始化每个玩家的 ID、手牌和状态
//	players := make([]*Player, playerCount)
//	for i := 0; i < playerCount; i++ {
//		players[i] = &Player{
//			PlayerID: playerId[i],
//			Cards:    make([]int, 0),
//			Status:   true,
//		}
//	}
//	// 创建一副新的牌，并洗牌
//	cardpool := NewDeck()
//	shuffle(cardpool)
//	// 将牌发给每个玩家，每人五张牌
//	for i := 0; i < playerCount*5; i++ {
//		card := cardpool[0]
//		cardpool = cardpool[1:]
//		players[i%playerCount].Cards = append(players[i%playerCount].Cards, card)
//	}
//	// 创建一个新的游戏实例，并返回指针
//	return &Game{
//		Players:   players,
//		CardPool:  cardpool,
//		Discarded: make([]int, 0),
//		TurnIndex: 0,
//		Order:     1,
//		IsOver:    true,
//	}
//}

// // 初始化牌堆
//
//	func NewDeck() []int {
//		cards := []int{}
//		cards = addCard(cards, 0, 8)
//		cards = addCard(cards, 1, 5)
//		cards = addCard(cards, 2, 3)
//		cards = addCard(cards, 3, 5)
//		cards = addCard(cards, 4, 3)
//		cards = addCard(cards, 5, 4)
//		cards = addCard(cards, 6, 3)
//		cards = addCard(cards, 7, 4)
//		cards = addCard(cards, 8, 4)
//		cards = addCard(cards, 9, 4)
//		cards = addCard(cards, 10, 4)
//		cards = addCard(cards, 11, 6)
//		return cards
//	}
//
// // 将牌按照数量加入牌组
// func addCard(cards []int, CardType int, nums int) []int {
//
//		for i := 0; i < nums; i++ {
//			cards = append(cards, CardType)
//		}
//		return cards
//	}
//
// // 玩家打出手中的某张牌
//
//	func (p *Player) PlayCard(game *Game, cardIndex int) {
//		log.Printf("玩家%s的回合\n", p.PlayerID)
//		// 从玩家手牌中取出要打出的卡牌
//		card := p.Cards[cardIndex]
//		// 将该卡牌从玩家手牌中删除，并将其加入弃牌堆中
//		p.Cards = append(p.Cards[:cardIndex], p.Cards[cardIndex+1:]...)
//		game.Discarded = append(game.Discarded, card)
//		card = 0
//		switch card {
//		case SKIP:
//			game.NextTurn()
//			return
//		case FORBID_1, FORBID_2:
//			game.blame(5)
//			return
//		case REVERSE:
//			game.reverse()
//			return
//		case DRAW_FROM_BOTTOM:
//			game.NextTurn()
//			p.drawbotton(game)
//			return
//		case DEMAND:
//			p.ask(game.Players[4], 0) //获取要选的卡的索引
//			return
//		case SWAP:
//			p.swap(game.Players[4], 0, 0)
//			return
//		case PREDICT:
//			predict := p.predict(game)
//			if predict == 13 {
//				log.Println("程序出错")
//			}
//			return
//		case VIEW:
//			p.see(game)
//			return
//		case SHUFFLE:
//			game.shuffleDeck()
//			return
//		case DISMANTLE:
//			//拆炸弹
//
//			return
//		default:
//			// 如果遇到未知卡牌类型，则报错提示
//			panic(fmt.Sprintf("invalid card type: %s", card))
//			return
//		}
//	}
//

//// 摸一张牌
//func (g *Game) draw() int {
//	card := g.CardPool[0]
//	g.CardPool = g.CardPool[1:]
//	return card
//}
//
//// 玩家抽牌
//func (p *Player) playerdraw(game *Game) int {
//	card := game.draw()
//
//	p.Cards = append(p.Cards, card)
//	return card
//}
//
//// 玩家使用洗牌
//func (g *Game) shuffleDeck() {
//	rand.Seed(time.Now().Unix())
//	for i := len(g.CardPool) - 1; i > 0; i-- {
//		j := rand.Intn(i + 1)
//		g.CardPool[i], g.CardPool[j] = g.CardPool[j], g.CardPool[i]
//	}
//}
//
//// 从底抽卡
//func (g *Game) drawbotton() int {
//	card := g.CardPool[len(g.CardPool)-1]
//	g.CardPool = g.CardPool[0 : len(g.CardPool)-2]
//	return card
//}
//
//// 反转函数
//func (g *Game) reverse() {
//	if g.Order == 1 { // 如果游戏顺序是正向
//		g.Order = -1
//		return
//	}
//	g.Order = 1
//}
//
//// 抽底
//func (p *Player) drawbotton(game *Game) {
//	card := game.drawbotton()
//	//抽到炸弹
//	if card == BOMB {
//
//	}
//	p.Cards = append(p.Cards, card)
//}
//
//// 转到下一个玩家回合
//func (g *Game) NextTurn() {
//	if g.Order == 1 { // 如果游戏顺序是正向
//		// 将轮到的玩家索引加一，如果超出玩家数量，则回到第一个玩家
//		g.TurnIndex = (g.TurnIndex + 1) % len(g.Players)
//	} else if g.Order == -1 { // 如果游戏顺序是反向
//		// 将轮到的玩家索引减一，如果小于零，则回到最后一个玩家
//		g.TurnIndex = (g.TurnIndex - 1 + len(g.Players)) % len(g.Players)
//	}
//
//}
//
//// 嫁祸卡
//func (g *Game) blame(targetIndex int) {
//	// 检查被嫁祸的玩家是否已经出局
//	if !g.Players[targetIndex].Status {
//		return
//	}
//	// 被嫁祸的玩家立即到该玩家的回合
//	g.TurnIndex = targetIndex
//	log.Printf("玩家%s的回合\n", g.Players[g.TurnIndex])
//
//}
//
//// 玩家索要另一个玩家的手牌
//func (p *Player) ask(play *Player, cardIndex int) {
//	// 判断对方的手牌是否为空，如果是则返回
//	if len(play.Cards) == 0 {
//		fmt.Printf("%s has no cards to give.\n")
//		return
//	}
//	card := play.Cards[cardIndex]                                            // 取出对方手牌中指定索引的卡牌
//	play.Cards = append(play.Cards[:cardIndex], play.Cards[cardIndex+1:]...) // 从对方手牌中移除被索要的卡牌
//	p.Cards = append(p.Cards, card)                                          // 将被索要的卡牌加入到自己的手牌中
//
//}
//
//// 预测下一张炸弹
//func (p *Player) predict(game *Game) int {
//	for i, card := range game.CardPool {
//		if card == BOMB {
//			return i
//		}
//	}
//	return 13 //返回13属于没有的类型
//}
//
//// 玩家使用透视牌，查看接下来三张牌的具体内容
//func (p *Player) see(game *Game) []int {
//	j := []int{}
//	if len(game.CardPool) < 3 { //少于三张直接透视全部
//		for i := 0; i < len(game.CardPool); i++ {
//			j[i] = game.CardPool[i]
//		}
//		return j
//	}
//	for i := 0; i < 3; i++ {
//		j = append(j, game.CardPool[i])
//	}
//	return j
//
//}
//
//// swap函数用于两个玩家交换手牌
//func (p *Player) swap(target *Player, cardIndex, mycardIndex int) {
//	// 如果目标玩家没有手牌，则直接返回
//	if len(target.Cards) == 0 {
//		return
//	}
//	// 取出需要交换的卡牌，并从自己手牌中移除
//	card := p.Cards[mycardIndex]
//	p.Cards = append(p.Cards[:mycardIndex], p.Cards[mycardIndex+1:]...)
//	// 取出需要接收的卡牌，并从目标玩家手牌中移除
//	targetCard := target.Cards[cardIndex]
//	target.Cards = append(target.Cards[:cardIndex], target.Cards[cardIndex+1:]...)
//	// 将取出的卡牌分别加入到另一个玩家的手牌中
//	p.Cards = append(p.Cards, targetCard)
//	target.Cards = append(target.Cards, card)
//
//}
//
//// 拆除炸弹
//func (p *Player) disarmBomb(g *Game) {
//	//拆除炸弹，炸弹回到牌堆
//	g.CardPool = append(g.CardPool, 10)
//	//加入后洗牌
//	g.shuffleDeck()
//	//然后到下一个玩家
//}
//
//// 游戏结束后输出结果
//func (g *Game) printResult() {
//	fmt.Println("Game over!")
//	for _, player := range g.Players {
//		status := "dead"
//		if player.Status {
//			status = "alive"
//		}
//		fmt.Println(status)
//	}
//}
//
//// 玩家出牌前30秒倒计时器
//func (p *Player) startTurnTimer(game *Game) {
//	timer := time.NewTimer(30 * time.Second)
//	defer timer.Stop()
//
//	select {
//	case <-timer.C:
//		// 时间到了还没有出牌，则自动出一张牌
//		p.PlayCard(game, 0)
//	case cardIndex := <-p.PlayCh:
//		// 玩家在规定时间内出牌
//		p.PlayCard(game, cardIndex)
//	}
//}
//
//// 玩家抽到炸弹牌时的拆弹倒计时器
//func (p *Player) startDefuseTimer() bool {
//	timer := time.NewTimer(10 * time.Second)
//	defer timer.Stop()
//
//	select {
//	case <-timer.C:
//		// 时间到了还没有选择，则默认不拆弹
//		return false
//	case choice := <-p.DefuseCh:
//		// 玩家做出了选择
//		return choice
//	}
//}
//
//// 判断有没有拆弹卡
//func (p *Player) hasDefuseCard() bool {
//	for _, card := range p.Cards {
//		if card == 11 {
//			return true
//		}
//	}
//	return false
//}
//func (g *Game) eliminate(p *Player) {
//	// 从游戏的玩家列表中删除该玩家
//	for i, player := range g.Players {
//		if player == p {
//			g.Players = append(g.Players[:i], g.Players[i+1:]...)
//			break
//		}
//	}
//
//}
//
//// 用于从用户获取布尔类型的输入。通常情况下，它会输出一些提示信息让用户选择“是”或“否”，然后返回相应的布尔值。
//func askForBooleanInput(prompt string) bool {
//	reader := bufio.NewReader(os.Stdin)
//	fmt.Print(prompt + " (y/n): ")
//	text, _ := reader.ReadString('\n')
//	text = strings.ToLower(strings.TrimSpace(text))
//	return text == "y" || text == "yes"
//}
//
//// 游戏循环，在每个玩家的回合执行
//func (game *Game) gameLoop(p *Player) {
//	p.startTurnTimer(game)
//	//抽卡
//	playerdraw := p.playerdraw(game)
//	// 如果玩家抽到炸弹牌，则询问是否拆弹
//	if playerdraw == 10 {
//		fmt.Println("You've got a bomb! Would you like to defuse it?")
//		// 向玩家发送一个拆弹请求，并开始倒计时
//		p.DefuseCh = make(chan bool)
//		go func() {
//			choice := askForBooleanInput("Y") //从客户端接收处理信息，要修改
//			p.DefuseCh <- choice
//		}()
//
//		// 启动拆弹倒计时器
//		isDefused := p.startDefuseTimer()
//
//		if !isDefused {
//			// 如果没有拆弹，尝试使用拆弹卡
//			if p.hasDefuseCard() {
//				p.disarmBomb(game)
//			} else {
//				// 没有拆弹卡，炸弹爆炸了，该玩家出局
//				game.eliminate(p)
//			}
//		}
//	}
//
//}
