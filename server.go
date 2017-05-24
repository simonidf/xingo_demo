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
	RunServer();
}

func RunServer(){
	go webserver.WebServerBase();

	network.NetWorkObj.AddRouter(&api.ApiRouter{});
	network.NetWorkObj.SetOnConnect(DoConnectionMade);
	network.NetWorkObj.SetOnClose(DoConnectionLost);
	network.NetWorkObj.Run();

	room.RoomMgrObj.Init();

	///logerr
}

func DoConnectionMade(fconn iface.Iconnection) {
	fmt.Printf("Connected")
	//p, _ := core.BattleFieldObj.AddPlayer(fconn)
	//fconn.SetProperty("pid", p.Pid)
}

func DoConnectionLost(fconn iface.Iconnection) {
	fmt.Printf("Lost");
	room.RoomMgrObj.OnPlayerLost(fconn);
}
