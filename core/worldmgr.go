package core

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/viphxin/xingo/iface"
	"xingo_demo/pb"
	"xingo_demo/ProtoTest"
	"sync"
	"github.com/viphxin/xingo/logger"
	//"time"
	"fmt"
	"time"
	"runtime"
)

//noinspection ALL
type WorldMgr struct {
	PlayerNumGen int32
	BulletNumGen int32
	Players      map[int32]*Player
	Bullets      map[int32]*Bullet
	AoiObj1       *AOIMgr//地图1
	MsgCollect *ProtoTest.FrameInfo
	sync.RWMutex
}

var WorldMgrObj *WorldMgr

func init() {
	WorldMgrObj = &WorldMgr{
		PlayerNumGen:    0,
		Players:         make(map[int32]*Player),
		BulletNumGen: 0,
		Bullets: make(map[int32]*Bullet),
		AoiObj1:          NewAOIMgr(85, 410, 75, 400, 10, 20),
		MsgCollect:&ProtoTest.FrameInfo{
			NewPlayer: make([]*ProtoTest.PlayerInfo, 0),
			PlayerMove: make([]*ProtoTest.PlayerInfo, 0),
			PlayerDead: make([]*ProtoTest.PlayerInfo, 0),
			AttackAction: make([]*ProtoTest.AttackAction, 0),
			Hit: make([]*ProtoTest.Hit, 0),
			SkillPrepare: make([]*ProtoTest.SkillInfo, 0),
			SkillCancel: make([]*ProtoTest.SkillInfo, 0),
			SkillStartUp:make([]*ProtoTest.SkillInfo, 0),
			PlayerRelife: make([]*ProtoTest.PlayerInfo, 0),
			ObjBorn: make([]*ProtoTest.ObjInfo, 0),
			ObjMove: make([]*ProtoTest.ObjInfo, 0),
			ObjDeleted: make([]*ProtoTest.ObjInfo, 0),
			SkillObjBorn: make([]*ProtoTest.SkillObjInfo, 0),
			SkillObjMove: make([]*ProtoTest.SkillObjInfo, 0),
			SkillObjDeleted: make([]*ProtoTest.SkillObjInfo, 0),
		},
	}

	go WorldMgrObj.LoopPush();
}

func (this *WorldMgr)AddPlayer(fconn iface.Iconnection) (*Player, error) {
	this.Lock()
	this.PlayerNumGen += 1
	p := NewPlayer(fconn, this.PlayerNumGen)
	this.Players[p.Pid] = p
	this.Unlock()
	//同步Pid
	msg := &pb.SyncPid{
		Pid: p.Pid,
	}
	p.SendMsg(1, msg)
	//加到aoi
	this.AoiObj1.Add2AOI(p)
	//周围的人
	p.SyncSurrouding()
	return p, nil
}

func (this *WorldMgr)AddBullet(pid int32) (*Bullet, error) {
	logger.Debug(fmt.Sprintf("AddBullet(%d)",pid));
	this.Lock()
	this.BulletNumGen += 1

	b := NewBullet(this.Players[pid],this,this.BulletNumGen)
	this.Bullets[b.Id] = b
	this.Unlock()

	logger.Debug(fmt.Sprintf("Generated(%d)",pid));

	//同步Pid
	msg := &ProtoTest.ObjInfo{
		Id : b.Id,
		X : b.X,
		Y :b.Y,
		Z :b.Z,
		V :b.V,
		Deleted:b.deleted,
	}

	this.MsgCollect.ObjBorn = append(this.MsgCollect.ObjBorn,msg);

	return b, nil
}

func (this *WorldMgr)RemovePlayer(pid int32){
	this.Lock()
	defer this.Unlock()
	//从aoi移除
	this.AoiObj1.LeaveAOI(this.Players[pid])
	delete(this.Players, pid)
}

func (this *WorldMgr)Move(p *Player){
	var data *pb.BroadCast
	data = &pb.BroadCast{
		Pid : p.Pid,
		Tp: 4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
			X: p.X,
			Y: p.Y,
			Z: p.Z,
			V: p.V,
			},
		},
	}
	/*aoi*/
	pids, err := this.AoiObj1.GetSurroundingPids(p)
	if err == nil{
		for _, pid := range pids{
			player, err1 := this.GetPlayer(pid)
			if err1 == nil{
				player.SendMsg(200, data)
			}
		}
	}
}

func (this *WorldMgr)SendMsgByPid(pid int32, msgId uint32, data proto.Message){
	p, err := this.GetPlayer(pid)
	if err == nil{
		p.SendMsg(msgId, data)
	}
}

func (this *WorldMgr) GetPlayer(pid int32)(*Player, error){
	this.RLock()
	defer this.RUnlock()
	p, ok := this.Players[pid]
	if ok{
		return p, nil
	}else{
		return nil, errors.New("no player in the world!!!")
	}
}

func (this *WorldMgr) Broadcast(msgId uint32, data proto.Message) {
	//logger.Debug(fmt.Sprintf("Broadcast: (%d)", msgId))
	this.RLock()
	defer this.RUnlock()
	for _, p := range this.Players {
		p.SendMsg(msgId, data)
	}
}

func (this *WorldMgr) BroadcastBuff(msgId uint32, data proto.Message) {
	this.RLock()
	defer this.RUnlock()
	for _, p := range this.Players {
		p.SendBuffMsg(msgId, data)
	}
}

func (this *WorldMgr) AOIBroadcast(p *Player, msgId uint32, data proto.Message) {
	/*aoi*/
	pids, err := WorldMgrObj.AoiObj1.GetSurroundingPids(p)
	if err == nil{
		for _, pid := range pids{
			player, err1 := WorldMgrObj.GetPlayer(pid)
			if err1 == nil {
				player.SendMsg(msgId, data)
			}
		}
	}else{
		logger.Error(err)
	}
}

func (this *WorldMgr) LoopPush() {
	go func() {
		for {
			//没有玩家后需要清除定时器
			if (true) {

				this.Step()
				//logger.Debug(fmt.Sprintf("BroadcastFrame"));
				this.Broadcast(301,this.MsgCollect);
				this.MsgCollect = &ProtoTest.FrameInfo{
					NewPlayer: make([]*ProtoTest.PlayerInfo, 0),
					PlayerMove: make([]*ProtoTest.PlayerInfo, 0),
					PlayerDead: make([]*ProtoTest.PlayerInfo, 0),
					AttackAction: make([]*ProtoTest.AttackAction, 0),
					Hit: make([]*ProtoTest.Hit, 0),
					SkillPrepare: make([]*ProtoTest.SkillInfo, 0),
					SkillCancel: make([]*ProtoTest.SkillInfo, 0),
					SkillStartUp:make([]*ProtoTest.SkillInfo, 0),
					PlayerRelife: make([]*ProtoTest.PlayerInfo, 0),
					ObjBorn: make([]*ProtoTest.ObjInfo, 0),
					ObjMove: make([]*ProtoTest.ObjInfo, 0),
					ObjDeleted: make([]*ProtoTest.ObjInfo, 0),
					SkillObjBorn: make([]*ProtoTest.SkillObjInfo, 0),
					SkillObjMove: make([]*ProtoTest.SkillObjInfo, 0),
					SkillObjDeleted: make([]*ProtoTest.SkillObjInfo, 0),
				};
				time.Sleep(time.Millisecond * time.Duration(30))
				runtime.GC()
			} else {
				//close(this.StepQueue)
				logger.Debug("LoopPush stoped successful!!!")
				return
			}
		}
	}()
}

func (this *WorldMgr) Step() {
	this.Update();
}

func (this *WorldMgr) Update() {
	for _, b := range this.Bullets {
		b.Update();
	}
}