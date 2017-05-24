package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"xingo_demo/room"
	"strings"
)

func WebServerBase() {
	fmt.Println("This is webserver base!")

	//第一个参数为客户端发起http请求时的接口名，第二个参数是一个func，负责处理这个请求。
	http.HandleFunc("/create_room_for_group", handleCreateRoom)

	//服务器要监听的主机地址和端口号
	err := http.ListenAndServe("127.0.0.1:30001", nil)

	if err != nil {
		fmt.Println("ListenAndServe error: ", err.Error())
	}
}

func handleCreateRoom(w http.ResponseWriter, req *http.Request) {
	fmt.Println("loginTask is running...")

	//模拟延时
	time.Sleep(time.Second * 1)

	//获取客户端通过GET/POST方式传递的参数
	req.ParseForm()
	credentials, found1 := req.Form["credentials"]

	fmt.Println(req)

	fmt.Println(credentials)


	if !(found1) {
		fmt.Fprint(w, "请勿非法访问")
		return
	}

	result := NewBaseJsonBean()

	var credentialArray []string

	credentialArray = strings.Split(credentials[0],",");

	if(room.RoomMgrObj.CreateNewRoomForPlayerGroup(credentialArray) == nil){
		result.Code = 200
		result.Message = "create room ok"
		result.Data = "";
		//向客户端返回JSON数据
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))
	}
}