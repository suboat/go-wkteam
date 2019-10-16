package wkteam

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// GetAgent 取开发者信息
func (api *WkTeam) GetAgent() (ret *Agent, err error) {
	var data = new(Agent)
	if _, err = api.Do("/foreign/user/getInfo", nil, data); err != nil {
		return
	}
	_ = data.init()
	ret = data
	return
}

// 发起请求
func (api *WkTeam) Do(name string, params map[string]interface{}, data interface{}) (ret []byte, err error) {
	if err = api.initApiKey(); err != nil {
		return
	}
	var (
		client  = &http.Client{}
		url     = fmt.Sprintf(`%s%s`, api.ApiHost, name)
		now     = time.Now().Unix()
		webTime = strconv.Itoa(int(now)) + "_s"
		token   = fmt.Sprintf(`%x`, md5.Sum([]byte(webTime+api.Secret)))
		msg     = &struct {
			Code int             `json:"code"` // 1成功，0失败
			Msg  string          `json:"msg"`  // 反馈信息
			Data json.RawMessage `json:"data"` //
		}{}
		req  *http.Request
		resp *http.Response
		raw  []byte
	)

	if req, err = http.NewRequest("POST", url, bytes.NewBuffer(nil)); err != nil {
		return
	}
	req.Header.Add("hswebtime", webTime)
	req.Header.Add("apikey", api.apiKey)
	req.Header.Add("token", token)
	if resp, err = client.Do(req); err != nil {
		return
	} else {
		raw, _ = ioutil.ReadAll(resp.Body)
		_ = resp.Close
	}
	if err = json.Unmarshal(raw, msg); err != nil {
		return
	}
	if msg.Code != 1 {
		// 失败
		if len(msg.Msg) > 0 {
			err = fmt.Errorf(`%s`, msg.Msg)
		} else {
			err = fmt.Errorf("unknown err: %s", string(raw))
		}
		return
	}

	// debug
	api.Log.Debugf(`[api-resp] %s -> %s`, name, string(raw))

	// 解析数据
	if data != nil {
		if err = json.Unmarshal(msg.Data, data); err != nil {
			return
		}
	}
	ret = msg.Data
	return
}

// 获取apiKey
func (api *WkTeam) initApiKey() (err error) {
	if len(api.apiKey) > 0 {
		return
	}
	api.lock.Lock()
	defer api.lock.Unlock()
	if len(api.apiKey) > 0 {
		return
	}
	// 请求获取key
	var (
		url  = fmt.Sprintf("%s/foreign/auth/login", api.ApiHost)
		resp *http.Response
		raw  []byte
		data = &struct {
			Code int    `json:"code"` // 1成功，0失败
			Msg  string `json:"msg"`  //
			Data struct {
				ApiKey string `json:"apikey"` // 有效期交互密钥
				Type   int    `json:"type"`   // 1管理员，12客服管理，13客户账号
			} `json:"data"`
		}{}
	)
	if resp, err = http.PostForm(url, map[string][]string{
		"phone":    {api.Phone},
		"password": {api.Password},
	}); err != nil {
		return
	} else {
		raw, _ = ioutil.ReadAll(resp.Body)
		_ = resp.Close
	}
	if err = json.Unmarshal(raw, data); err != nil {
		return
	} else if data.Code != 1 {
		return fmt.Errorf(`%s`, data.Msg)
	}
	api.apiKey = data.Data.ApiKey

	// debug
	api.Log.Debugf(`[api-init] get apiKey:%s`, api.apiKey)

	//
	return
}
