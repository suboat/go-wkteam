package wkteam

import (
	//"github.com/suboat/go-contrib"

	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/suboat/go-contrib"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

// GetGroups 取群列表
func (api *WkTeam) GetGroups(query *Query) (ret []*Group, err error) {
	if query == nil {
		query = &Query{}
	}
	if len(query.Account) == 0 {
		if query.Account = api.Account; len(query.Account) == 0 {
			query.Account = Settings.Account
		}
	}
	var (
		data []*Group
	)
	if _, err = api.Do("/foreign/group/get", query, &data); err != nil {
		return
	}
	ret = data
	return
}

// Do 发起请求
func (api *WkTeam) Do(name string, query *Query, data interface{}) (ret []byte, err error) {
	if err = api.initApiKey(); err != nil {
		return
	}
	var (
		client  = &http.Client{}
		uri     = fmt.Sprintf(`%s%s`, api.ApiHost, name)
		now     = time.Now().Unix()
		webTime = strconv.Itoa(int(now)) + "_s"
		token   = fmt.Sprintf(`%x`, md5.Sum([]byte(webTime+api.Secret)))
		msg     = &struct {
			Code int             `json:"code"` // 1成功，0失败
			Msg  string          `json:"msg"`  // 反馈信息
			Data json.RawMessage `json:"data"` //
		}{}
		//
		account = ""
		limit   = 10
		skip    = 0
		params  = map[string][]string{}
		req     *http.Request
		resp    *http.Response
		raw     []byte
	)

	// 参数
	if query != nil {
		if len(query.Account) > 0 {
			account = query.Account
		}
		if query.Limit > 0 {
			if limit = query.Limit; limit > 100 {
				limit = 100 // limit最大值100
			}
		}
		if query.Skip > 0 {
			skip = query.Skip
		}
		if query.Page > 1 {
			// page1起始
			skip = (query.Page - 1) * limit
		}
		if query.Params != nil {
			for k, v := range query.Params {
				params[k] = []string{fmt.Sprintf("%v", v)}
			}
		}
		if len(account) > 0 {
			params["account"] = []string{account}
		}
		if limit > 0 {
			if skip > 0 {
				if _page := skip/limit + 1; _page > 0 {
					params["num"] = []string{fmt.Sprintf("%d", limit)}
					params["page"] = []string{fmt.Sprintf("%d", _page)}
				}
			}
		}
	}
	if req, err = http.NewRequest("POST", uri, strings.NewReader(url.Values(params).Encode())); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
			err = fmt.Errorf(`%s <- %s`, msg.Msg, PubJSON(params))
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

// PassAddFriend 同意添加好友
func (api *WkTeam) PassAddFriend(account string) (ret bool, err error) {
	if len(account) == 0 {
		// 微信号必填
		err = contrib.ErrParamUndefined
		return
	}
	var (
		param = map[string]interface{}{
			"my_account": Settings.Account,
			"account":    account, // 好友微信号
		}
	)
	if _, _err := api.Do("/foreign/friends/passAddFriends", &Query{Params: param}, nil); _err != nil {
		err = _err
		return
	} else {
		// 成功添加朋友
		ret = true
	}
	return
}

// RemarkFriend 备注好友
func (api *WkTeam) RemarkFriend(account, remarkName string) (ret bool, err error) {
	if len(account) == 0 {
		// 微信号必填
		err = contrib.ErrParamUndefined
		return
	}
	var (
		param = map[string]interface{}{
			"my_account": Settings.Account,
			"to_account": account,    // 好友微信号
			"remark":     remarkName, // 备注名
		}
	)
	if _, _err := api.Do("/foreign/friends/remark", &Query{Params: param}, nil); _err != nil {
		err = _err
		return
	} else {
		// 成功备注好友
		ret = true
	}
	return
}
