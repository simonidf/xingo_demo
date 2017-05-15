package main

import (
	"github.com/viphxin/xingo/iface"
	"xingo_demo/api"
	"xingo_demo/core"
	"xingo_demo/network"
	_ "net/http"
	_ "net/http/pprof"
	_ "runtime/pprof"
	_ "time"
	"fmt"
)

func main() {
	network.NetWorkObj.AddRouter(&api.ApiRouter{});
	network.NetWorkObj.SetOnConnect(DoConnectionMade);
	network.NetWorkObj.SetOnClose(DoConnectionLost);
	network.NetWorkObj.Run();

	core.BattleFieldObj.Init();
}

func DoConnectionMade(fconn iface.Iconnection) {
	fmt.Printf("Connected")
	p, _ := core.BattleFieldObj.AddPlayer(fconn)
	fconn.SetProperty("pid", p.Pid)
}

func DoConnectionLost(fconn iface.Iconnection) {
	fmt.Printf("Lost")
	pid, _ := fconn.GetProperty("pid")
	p, _ := core.BattleFieldObj.GetPlayer(pid.(int32))
	//移除玩家
	core.BattleFieldObj.RemovePlayer(pid.(int32))
	//消失在地图
	p.LostConnection()
}
