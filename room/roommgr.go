package room

import (
	"github.com/viphxin/xingo/iface"
	"xingo_demo/core"
	"sync"
)

type RoomMgr struct {
	RoomNumGen int32
	BattleFields      map[int32]*core.BattleField
	PidToBattleId	map[int32]int32
	sync.RWMutex
}

var RoomMgrObj *RoomMgr

func init() {
	RoomMgrObj = &RoomMgr{
		RoomNumGen:    0,
		BattleFields:         make(map[int32]*core.BattleField),
		PidToBattleId:         make(map[int32]int32),
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

func (this *RoomMgr)GetBidForPid(pid int32) (int, error) {

	//bid = this.PidToBattleId[pid];

	return 1, nil
}

func (this *RoomMgr)CreateNewRoomForPlayerGroup() (*core.BattleField, error) {
	this.Lock();
	this.RoomNumGen += 1;
	battleField := core.NewBattleField();
	this.BattleFields[this.RoomNumGen] = battleField;
	this.Unlock();
	return battleField, nil
}

func (this *RoomMgr)AddPlayerToBattle(fconn iface.Iconnection,bid int32){

}
