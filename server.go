package main

import (
	"github.com/viphxin/xingo/iface"
	"xingo_demo/api"
	"xingo_demo/core"
	"xingo_demo/network"
	"xingo_demo/room"
	_ "net/http"
	_ "net/http/pprof"
	_ "runtime/pprof"
	_ "time"
	"fmt"
	"io"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world hahahah!\n")


}

func HttpServer(){
http.HandleFunc("/hello", HelloServer)
err := http.ListenAndServe(":12345", nil)
if err != nil {
log.Fatal("ListenAndServe: ", err)
}else{
fmt.Printf("ListenAndServe:12345");
}
}

func main() {
	go HttpServer();

	network.NetWorkObj.AddRouter(&api.ApiRouter{});
	network.NetWorkObj.SetOnConnect(DoConnectionMade);
	network.NetWorkObj.SetOnClose(DoConnectionLost);
	network.NetWorkObj.Run();

	core.BattleFieldObj.Init();

	room.RoomMgrObj.Init();


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
