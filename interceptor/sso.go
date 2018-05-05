package interceptor

import (
	"net/http"
	"github.com/scofier/gogateway/util"
	"errors"
)

func SsoHandler()  {
	http.HandleFunc("/sso/login", userLogin)
	http.HandleFunc("/sso/checkToken", checkToken)
	http.HandleFunc("/sso/refreshToken", refreshToken)
}

func userLogin (w http.ResponseWriter, req *http.Request)  {
	var token string
	var err error
	var username = req.FormValue("username")
	var password = req.FormValue("password")


	if len(username) > 1 && len(password) > 1 {
		token =  util.GenerateRandnumString()
		err = util.SetRedisValue(token, username, -1)
	}else{
		err = errors.New("Login info error !")
	}
	util.JSONResp(w, err, token)
}

func checkToken (w http.ResponseWriter, req *http.Request)  {
	var token = req.Header.Get(util.HEAD_TOKEN)
	value,err := util.GetRedisValue(token)

	util.JSONResp(w, err, value)
}

func refreshToken (w http.ResponseWriter, req *http.Request)  {
	var token = req.Header.Get(util.HEAD_TOKEN)
	var value,err = util.GetRedisValue(token)
	if nil == err {
		token = util.GenerateRandnumString()
		err = util.SetRedisValue(token, value, -1)
	}
	util.JSONResp(w, err, token)
}