package main

import (
	"paas/icsoc-monitor/config"

	"context"
	"time"

	"github.com/alecthomas/log4go"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/client/http"
	"github.com/micro/go-plugins/registry/etcdv3"
)

const(
	SERVER_NAME              = "IM-CS"
	METHOD_NAME 			 = "/im-cs/syncToMonitor"
)

func main(){
	config.Init("")
	conf := config.GetConf()
	reg := etcdv3.NewRegistry(registry.Addrs(conf.Etcd.Addr))
	cli := http.NewClient(client.Registry(reg))

	type JsonRequest struct {
		Params string `json:"params"`
		ItemId string `json:"item_id"`
	}

	type JsonResponse struct {
		Code int `json:"code"`
		Message string `json:"message"`
	}

	reqS := JsonRequest{}
	rspS := JsonResponse{}
	req := cli.NewJsonRequest(SERVER_NAME, METHOD_NAME, &reqS)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
	callErr := cli.Call(ctx, req, &rspS)
	cancel()

	if callErr != nil {
		log4go.Info("call fail..err:%v", callErr)
	}else {
		log4go.Info("call suc.....")
	}

	time.Sleep(time.Second * 2)
}
