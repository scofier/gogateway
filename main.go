package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/scofier/gogateway/proxy"
	"encoding/json"
	"github.com/scofier/gogateway/interceptor"
	"github.com/scofier/gogateway/util"
	"flag"
)

var port=":9090"

func main() {
	var configFile = flag.String("c", "config.json", "Config file, Usage: -c config.json")
	flag.Parse()
	//加载配置
	var config proxy.DefaultRegistry
	util.Load(*configFile, &config)
	//启动反向代理
	http.HandleFunc("/", proxy.NewMultipleHostReverseProxy(config, nil))
	//状态信息
	http.HandleFunc("/info", func(w http.ResponseWriter, req *http.Request) {
		value, _:=json.Marshal(config)
		fmt.Fprintf(w, "%s", value)
	})
	//加载sso模块
	interceptor.SsoHandler()
	//启动服务
	println("Start on", port)
	log.Fatal(http.ListenAndServe(port, nil))
}