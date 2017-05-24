package api

import (
	"xingo_demo/pb"
	"xingo_demo/ProtoTest"
	"xingo_demo/core"
	"xingo_demo/room"
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

			cre, _ := request.Fconn.GetProperty("cre");

			room := room.RoomMgrObj.GetRoomByCre(cre.(string));

			p, _ := room.GetPlayer(pid.(int32))
			p.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)

			room.Move(p);
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
	fmt.Printf("Api_4")
	msg := &ProtoTest.AttackAction{}
	err := proto.Unmarshal(request.Pdata.Data, msg)
	if err == nil {
		pid, err1 := request.Fconn.GetProperty("pid")
		if err1 == nil{
			logger.Debug(fmt.Sprintf("AttackAction: (%d)", pid))

			cre, _ := request.Fconn.GetProperty("cre");

			room := room.RoomMgrObj.GetRoomByCre(cre.(string));

			p, _ := room.GetPlayer(pid.(int32))
			room.AddBullet(p.Pid);
		}else{
			logger.Error(err1)
			request.Fconn.LostConnection()
		}



	} else {
		logger.Error(err)
		request.Fconn.LostConnection()
	}
}

func (this *ApiRouter) Api_5(request *fnet.PkgAll) {
	fmt.Printf("Api_5")
	msg := &ProtoTest.CredentialInfo{}
	err := proto.Unmarshal(request.Pdata.Data, msg)
	if err == nil {
		logger.Debug(fmt.Sprintf("CredentialInfo: (%s)", msg.Cre))

		bid:=room.RoomMgrObj.GetBidForCre(msg.Cre);

		fmt.Println("bid:");
		fmt.Println(bid);

		request.Fconn.SetProperty("cre",msg.Cre);

		fmt.Println("2222");

		room.RoomMgrObj.AddPlayerToBattle(request.Fconn,bid);

		fmt.Println("3333");

	} else {
		logger.Error(err)
		request.Fconn.LostConnection()
	}
}