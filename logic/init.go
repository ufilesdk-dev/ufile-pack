package logic

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	uflog "ufile-pack/gosdk/log"
)

const (
	GetUFileZip          = "GetUFileZipRequest"
	GetUFileZipByList    = "GetUFileZipByListRequest"
	GetUFileZipByListExt = "GetUFileZipByListExtRequest"
)

func HttpRouter(w http.ResponseWriter, r *http.Request) {
	var response []byte
	var err error

	body, _ := ioutil.ReadAll(r.Body)
	var tmp map[string]interface{}
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		//uflog.ERROR("HttpRouter|unmarshal err:", err)
		return
	}

	if _, ok := tmp["action"]; !ok {
		uflog.ERROR("HttpRouter| no action | req :", tmp)
		return
	}

	switch tmp["action"] {
	case GetUFileZip:
		response, err = GetZipFileRequest(body)
	case GetUFileZipByList:
		response, err = GetZipFileByListRequest(body)
	case GetUFileZipByListExt:
		response, err = GetZipFileByListExtRequest(body)
	default:
		response = []byte("Action Not Support！ ")
	}

	if err != nil {
		uflog.ERROR("HttpRouter|handle action|err:", err)
		return
	}

	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Length", strconv.Itoa(len(response)))
	w.Write(response)
}
