package network

import (
"github.com/viphxin/xingo/fserver"
"github.com/viphxin/xingo/iface"
"github.com/viphxin/xingo/utils"
)

//noinspection ALL
type NetWork struct {
	OnConnectioned         func(fconn iface.Iconnection)
	OnClosed               func(fconn iface.Iconnection)
	netserver		iface.Iserver
	router			interface{}
}

var NetWorkObj *NetWork

func init() {
	NetWorkObj = &NetWork{
		netserver:fserver.NewServer(),
	}
}

func (this *NetWork) SetOnConnect(onConnect func(fconn iface.Iconnection)) {
	utils.GlobalObject.OnConnectioned = onConnect
}

func (this *NetWork) SetOnClose(onClose func(fconn iface.Iconnection)) {
	utils.GlobalObject.OnClosed = onClose;
}

func (this *NetWork) AddRouter(_router interface{}) {
	this.router = _router;
	this.netserver.AddRouter(this.router);
}

func (this *NetWork) Run(){
	this.netserver.Serve()
}

