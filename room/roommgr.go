package room

import (
	"github.com/viphxin/xingo/iface"
	"xingo_demo/core"
	"sync"
	"fmt"
)

type RoomMgr struct {
	RoomNumGen int32
	BattleFields      map[int32]*core.BattleField
	CreToBattleId	map[string]int32
	PidToBattleId map[int32]int32
	sync.RWMutex
}

var RoomMgrObj *RoomMgr

func init() {
	RoomMgrObj = &RoomMgr{
		RoomNumGen:    0,
		BattleFields:         make(map[int32]*core.BattleField),
		PidToBattleId:         make(map[int32]int32),
		CreToBattleId:  make(map[string]int32),
	}
}

func (this *RoomMgr)Init(){

}


func (this *RoomMgr)CreateNewRoom() (*core.BattleField, error) {
	this.Lock();
	this.RoomNumGen += 1;
	battleField := core.NewBattleField();
	this.BattleFields[this.RoomNumGen] = battleField;
	this.Unlock();
	return battleField, nil
}

func (this *RoomMgr)GetBidForCre(cre string) (int32) {

	bid := this.CreToBattleId[cre];

	return bid
}

func (this *RoomMgr)CreateNewRoomForPlayerGroup(credentialArray []string) (error) {
	this.Lock();

	battleField := core.NewBattleField();

	fmt.Println(battleField);

	battleField.Bid = this.RoomNumGen;

	for _, cre := range credentialArray{
		fmt.Println(cre);
		this.CreToBattleId[cre] = battleField.Bid;
		fmt.Println(battleField.Bid)
	}

	this.BattleFields[this.RoomNumGen] = battleField;

	this.RoomNumGen += 1;
	this.Unlock();
	return  nil;
}

func (this *RoomMgr)AddPlayerToBattle(fconn iface.Iconnection,bid int32){

	fmt.Println("AddPlayerToBattle bid " + string(bid));

	fmt.Println("AddPlayerToBattle1111");

	room := this.BattleFields[bid];

	fmt.Println("AddPlayerToBattle2222")
	fmt.Println(room);
	p, _ := room.AddPlayer(fconn);

	fconn.SetProperty("pid", p.Pid);

	fmt.Println("AddPlayerToBattle3333")


	if(len(room.Players) == 1){
		room.RunFrameRate();
	}
}

func (this *RoomMgr)GetBidByCre(cre string) int32 {
	fmt.Printf("GetBidByCre");

	bid := this.CreToBattleId[string(cre)];

	return bid;
}

func (this *RoomMgr)GetRoomByCre(cre string) *core.BattleField {
	fmt.Printf("GetBidByCre");

	bid := this.CreToBattleId[string(cre)];

	return this.BattleFields[bid];
}

func (this *RoomMgr)OnPlayerLost(fconn iface.Iconnection) {
	fmt.Printf("OnPlayerLost");

	pid, _ := fconn.GetProperty("pid")
	cre, _ := fconn.GetProperty("cre")

	bid := this.GetBidByCre(cre.(string));

	battleRoom := this.BattleFields[bid];

	p, _ := battleRoom.GetPlayer(pid.(int32))
	//移除玩家
	battleRoom.RemovePlayer(pid.(int32))
	//消失在地图
	p.LostConnection();
}
