package core

import (
	"math"
	"fmt"
	"github.com/viphxin/xingo/logger"
	"xingo_demo/ProtoTest"
	"time"
)

type Bullet struct {
	FromP   *Player
	battleField *BattleField
	AD int32
	Id int32
	X     float32//平面x
	Y     float32//高度
	Z     float32//平面y!!!!!注意不是Y
	V     float32//旋转0-360度
	Speed float32
	deleted bool
	lifetime int64
	starttime int64
}

func NewBullet(fromPlayer *Player,_battleField *BattleField,bid int32) *Bullet {

	logger.Debug(fmt.Sprintf("StartNewBullet"));

	logger.Debug(fmt.Sprintf("NewBullet"));

	b := &Bullet{
		battleField: _battleField,
		FromP:   fromPlayer,
		AD: fromPlayer.AD,
		Id: bid,
		X:     fromPlayer.X,
		Y:    fromPlayer.Y + 1,
		Z:     fromPlayer.Z,
		V:     fromPlayer.V,
		Speed: 0.5,
		lifetime:0,
		starttime: time.Now().Unix(),
	}

	logger.Debug(fmt.Sprintf("NewBulletEnd"));
	return b
}

func (this *Bullet) Update() {
	if(this.deleted){return};

	//更新
	xDelta := math.Sin(float64(this.V/180*3.14)) * float64(this.Speed);
	zDelta := math.Cos(float64(this.V/180*3.14)) * float64(this.Speed);
	this.X+=float32(xDelta);
	this.Z+=float32(zDelta);

	this.CheckHit();

	msg := &ProtoTest.ObjInfo{
		Id : this.Id,
		X : this.X,
		Y :this.Y,
		Z :this.Z,
		V :this.V,
		Deleted:this.deleted,
	}

	//logger.Debug(fmt.Sprintf("Broadcast(%d)",102));

	this.lifetime = time.Now().Unix() - this.starttime;

	//fmt.Printf("lifetime:(%d)",this.lifetime);

	if(this.lifetime>=2){
		this.deleted = true;
		msgDeleted := &ProtoTest.ObjInfo{
			Id : this.Id,
			X : this.X,
			Y :this.Y,
			Z :this.Z,
			V :this.V,
			Deleted:this.deleted,
		}
		this.battleField.MsgCollect.ObjDeleted = append(this.battleField.MsgCollect.ObjDeleted,msgDeleted);
	}

	this.battleField.MsgCollect.ObjMove = append(this.battleField.MsgCollect.ObjMove,msg);
}

func (this *Bullet) CheckHit() {
	for pid, player := range this.battleField.Players {

		if((player.Pid!=this.FromP.Pid) == true){
			if(this.CheckPlayer(player)){
				var temp = pid;
				temp = temp;
				this.deleted = true;

				msg := &ProtoTest.ObjInfo{
					Id : this.Id,
					X : this.X,
					Y :this.Y,
					Z :this.Z,
					V :this.V,
					Deleted:this.deleted,
				}
				this.battleField.MsgCollect.ObjMove = append(this.battleField.MsgCollect.ObjMove,msg);

				msg2 := &ProtoTest.Hit{
					Pid:player.Pid,
					HitHp:10,
				}
				this.battleField.MsgCollect.Hit = append(this.battleField.MsgCollect.Hit,msg2);

				break;
			}
		}

	}
}

func (this *Bullet) CheckPlayer(p *Player) bool{
	var distance = float64((p.X - this.X) * (p.X - this.X) + (p.Z - this.Z) * (p.Z - this.Z));
	distance = math.Pow(distance,0.5);

	if(distance<0.5){
		return true;
	}
	return false;
}
