package interceptor

import (
	"net/http"
	"github.com/scofier/gogateway/util"
	"errors"
)

func JwtInterceptor(w http.ResponseWriter, req *http.Request) error {

	var token = req.Header.Get(util.HEAD_TOKEN)
	err := errors.New("Access Forbidden!")
	if len(token) > 1 {
		_,err = util.GetRedisValue(token)
		if nil != err {
			err = errors.New("Token error!")
		}
	}
	if nil != err {
		util.JSONResp(w, err, nil)
	}
	return err
}

