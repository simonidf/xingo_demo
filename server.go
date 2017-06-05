package main

import (
	"github.com/viphxin/xingo/iface"
	"xingo_demo/api"
	"xingo_demo/network"
	"xingo_demo/room"
	"xingo_demo/webserver"
	_ "net/http"
	_ "net/http/pprof"
	_ "runtime/pprof"
	_ "time"
	"fmt"
)

func main() {
	RunServer();
}

func RunServer(){
	go webserver.WebServerBase();

	network.NetWorkObj.AddRouter(&api.ApiRouter{});
	network.NetWorkObj.SetOnConnect(DoConnectionMade);
	network.NetWorkObj.SetOnClose(DoConnectionLost);
	network.NetWorkObj.Run();

	room.RoomMgrObj.Init();

	fmt.Printf("Running...");
}

func DoConnectionMade(fconn iface.Iconnection) {
	fmt.Printf("Connected")
	//do nothing but wait for player send cre to join a room
}

func DoConnectionLost(fconn iface.Iconnection) {
	fmt.Printf("Lost");
	room.RoomMgrObj.OnPlayerLost(fconn);
}
