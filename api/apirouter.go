package api

import (
	"xingo_demo/pb"
	"xingo_demo/ProtoTest"
	"xingo_demo/core"
	"github.com/golang/protobuf/proto"
	"github.com/viphxin/xingo/fnet"
	"github.com/viphxin/xingo/logger"
	_ "time"
	"fmt"
	"github.com/viphxin/xingo/utils"
)

type ApiRouter struct {
}

/*
ping test
*/
func (this *ApiRouter) Api_0(request *fnet.PkgAll) {
	logger.Debug("call Api_0")
	// request.Fconn.SendBuff(0, nil)
	packdata, err := utils.GlobalObject.Protoc.GetDataPack().Pack(0, nil)
	if err == nil{
		request.Fconn.Send(packdata)
	}else{
		logger.Error("pack data error")
	}
}

/*
世界聊天
 */
func (this *ApiRouter) Api_2(request *fnet.PkgAll) {
	msg := &pb.Talk{}
	err := proto.Unmarshal(request.Pdata.Data, msg)
	if err == nil {
		logger.Debug(fmt.Sprintf("user talk: content: %s.", msg.Content))
		pid, err1 := request.Fconn.GetProperty("pid")
		if err1 == nil{
			p, _ := core.BattleFieldObj.GetPlayer(pid.(int32))
			p.Talk(msg.Content)
		}else{
			logger.Error(err1)
			request.Fconn.LostConnection()
		}

	} else {
		logger.Error(err)
		request.Fconn.LostConnection()
	}
}

/*
移动
 */
func (this *ApiRouter) Api_3(request *fnet.PkgAll) {
	msg := &pb.Position{}
	err := proto.Unmarshal(request.Pdata.Data, msg)
	if err == nil {
		logger.Debug(fmt.Sprintf("user move: (%f, %f, %f, %f)", msg.X, msg.Y, msg.Z, msg.V))
		pid, err1 := request.Fconn.GetProperty("pid")
		if err1 == nil{
			p, _ := core.BattleFieldObj.GetPlayer(pid.(int32))
			p.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
		}else{
			logger.Error(err1)
			request.Fconn.LostConnection()
		}

	} else {
		logger.Error(err)
		request.Fconn.LostConnection()
	}
}

/*
hahahah
 */
func (this *ApiRouter) Api_4(request *fnet.PkgAll) {
	msg := &ProtoTest.AttackAction{}
	err := proto.Unmarshal(request.Pdata.Data, msg)
	if err == nil {
		pid, err1 := request.Fconn.GetProperty("pid")
		if err1 == nil{
			logger.Debug(fmt.Sprintf("AttackAction: (%d)", pid))

			p, _ := core.BattleFieldObj.GetPlayer(pid.(int32))
			core.BattleFieldObj.AddBullet(p.Pid);
		}else{
			logger.Error(err1)
			request.Fconn.LostConnection()
		}



	} else {
		logger.Error(err)
		request.Fconn.LostConnection()
	}
}