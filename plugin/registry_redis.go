package plugin


import (
	"net/http"
	"github.com/scofier/gogateway/proxy"
	"fmt"
	"log"
	"strings"
	"time"
	"net"
	"github.com/scofier/gogateway/util"
)


// 在main.go 中增加下面的配置即可启用
// manage.RegistryHandler(&serviceRegistry)

func RegistryHandler(reg *proxy.DefaultRegistry)  {
	//定时同步缓存数据
	//go scheduleRunSync(reg)

	http.HandleFunc("/manage/addServer", addServer(reg))
	http.HandleFunc("/manage/delServer", delServer(reg))
	http.HandleFunc("/manage/serverInfo", serverInfo(reg))
	http.HandleFunc("/manage/flush", func(http.ResponseWriter, *http.Request) {
		syncServerToLocal(reg)
	})
}

func addServer(reg *proxy.DefaultRegistry) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		name,version,server,err := util.GetRequestInfo(req)
		if nil == err {
			value:=fmt.Sprintf("%s/%s/%s", name, version, server)

			if !util.CheckRedisSetValue(util.SERVER_KYE_IN_REDIS, value) {
				util.AddRedisSetValuePerm(util.SERVER_KYE_IN_REDIS, value)
				reg.Add(name,version,server)
			}
		}
		util.JSONResp(w, err, nil)
	}
}


func delServer(reg *proxy.DefaultRegistry) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {

		name,version,server,err := util.GetRequestInfo(req)
		if nil == err {
			reg.Delete(name,version,server)
		}
		util.JSONResp(w, err, nil)
	}
}


func serverInfo(reg *proxy.DefaultRegistry) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		servers := util.GetRedisSetList(util.SERVER_KYE_IN_REDIS)
		util.JSONResp(w, nil, servers)
	}
}

// schedule work
func scheduleRunSync(reg *proxy.DefaultRegistry) {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-t.C:
			syncServerToLocal(reg)
		}
	}
}

func syncServerToLocal(reg *proxy.DefaultRegistry)  {
	localValues :=reg.Values()
	servers := util.GetRedisSetList(util.SERVER_KYE_IN_REDIS)
	//缓存为空不更新,数据相等不更新
	if len(servers) == 0 || checkEQS(localValues, servers){
		return
	}else{
		log.Printf("[WARN] synchronized start >>>>\nL:%v\nS:%v", localValues, servers)
	}
	var modified = false
	//add new server
	for _,server := range servers{
		var servIsNew = true
		for sv,_ := range localValues {
			if sv == server {
				servIsNew = false
			}
		}
		if servIsNew {
			v := strings.Split(server, "/")
			modified = true
			if checkConnectOK(v[2]) {
				reg.Add(v[0], v[1], v[2])
			}
		}
	}
	//remove server
	for local,_ := range localValues {
		var localNeedDel = true
		for _,sv := range servers {
			if sv == local {
				localNeedDel = false
			}
		}
		if localNeedDel {
			v := strings.Split(local, "/")
			modified = true
			reg.Delete(v[0], v[1], v[2])
		}
	}
	if modified {
		log.Printf("[WARN] synchronized ended <<<<")
	}
}

func checkEQS(a map[string]string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _,node := range b {
		a[node] = node
	}
	if len(a) == len(b) {
		return true
	} else {

		return false
	}
}

func checkConnectOK(server string) bool {
	_, err := net.Dial("tcp", server)
	if err != nil {
		log.Printf("[WARN] server connect failed: [%v]", server)
		return false
	} else {
		return true
	}
}