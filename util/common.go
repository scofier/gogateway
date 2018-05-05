package util

import (
	"net/http"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"math/rand"
	io "io/ioutil"
)


const SERVER_KYE_IN_REDIS = "_SERVER_KEY_IN_REDIS"

const (
	//需要一个新的redis
	REDIS_host  = "127.0.0.1:6379"
	REDIS_pswd  = ""
	REDIS_dbNum = 0
)


const (
	HEAD_TOKEN = "X-Token"
)

type Result struct {
	ResultCode interface{} `json:"resultCode"`
	ResultMsg  interface{} `json:"resultMessage"`
	ResultData interface{} `json:"resultData"`
}




func GenerateRandnumString() string {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))

	return vcode
}


func JSONResp(w http.ResponseWriter, status error, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var statusCode = 0
	errs := ""
	if nil != status {
		statusCode = -1
		errs = status.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}else{
		w.WriteHeader(http.StatusOK)
	}

	result := Result{
		ResultCode: statusCode,
		ResultMsg:  errs,
		ResultData: v,
	}

	var serr = json.NewEncoder(w).Encode(result)
	if serr != nil {
		panic(serr)
	}
}

func GetRequestInfo(req *http.Request) (name,version,server string, err error) {
	name=req.FormValue("name")
	version=req.FormValue("version")
	server=req.FormValue("server")

	if len(name) < 1 || len(version) < 1 || len(server) < 1 {
		err = errors.New("request parameter error!")
	}
	return
}


func Load (filename string, v interface{}) bool {
	data, err := io.ReadFile(filename)
	if err != nil {
		fmt.Println("Load config json failed ===> filename : "+filename+" -- "+err.Error())
		return false
	}
	datajson := []byte(data)
	err = json.Unmarshal(datajson, v)
	if err != nil{
		fmt.Println("Read json failed ===> filename : "+filename+" -- "+err.Error())
		return false
	}
	return true
}