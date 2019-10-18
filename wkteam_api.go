package wkteam

import (
	"github.com/suboat/go-contrib"

	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

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
			Code  int             `json:"code"`  // 1成功，0失败
			Msg   string          `json:"msg"`   // 反馈信息
			Total int             `json:"total"` // 总记录数
			Data  json.RawMessage `json:"data"`  //
		}{}
		//
		account = ""
		limit   = 0
		skip    = 0
		params  = map[string][]string{}
		req     *http.Request
		resp    *http.Response
		raw     []byte
	)

	// 参数
	if query != nil {
		// 更新参数
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
		if query.Page > 0 && limit > 0 {
			skip = (query.Page) * limit
		}
		if query.Params != nil {
			for k, v := range query.Params {
				params[k] = []string{fmt.Sprintf("%v", v)}
			}
		}
		// 写入参数
		if len(account) > 0 && params["account"] == nil {
			params["account"] = []string{account}
		}
		if limit > 0 {
			params["num"] = []string{fmt.Sprintf("%d", limit)}
			if skip > 0 {
				// page是1起始
				if _page := skip/limit + 1; _page > 0 {
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

	// 反馈总记录数
	if msg.Total > 0 && query != nil {
		query.Total = msg.Total
	}

	// debug
	api.Log.Debugf(`[api-resp] %s %s -> %s`, name, PubJSON(params), string(raw))

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
	if err = api.init(); err != nil {
		return
	} else if len(api.apiKey) > 0 {
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

// 登记回调
func (api *WkTeam) SetCallback(urlPrefix string) (err error) {
	urlPrefix = strings.TrimRight(urlPrefix, "/")
	var (
		params = map[string]interface{}{
			"callbackSend": urlPrefix + CallbackPrefix + CallbackNormal,   // 通用回调
			"crowdlog":     urlPrefix + CallbackPrefix + CallbackMsgGroup, // 群聊回调
		}
	)
	if _, err = api.Do("/foreign/user/setUrl", &Query{Params: params}, nil); err != nil {
		return
	}
	// debug
	api.Log.Infof(`[api-config] callback params: %s`, PubJSON(params))
	return
}

// GetAgent 取开发者信息
func (api *WkTeam) GetAgent() (ret *Agent, err error) {
	var (
		data = &struct {
			ID          int    `json:"uid"`                    // 开发者uid
			Phone       string `json:"phone"`                  // 开发者手机号
			Name        string `json:"name"`                   // 开发者名称
			Sex         string `json:"sex,omitempty"`          // 性别
			OverdueTime int64  `json:"overdue_time,omitempty"` //
			LastTime    int64  `json:"last_time,omitempty"`    //
		}{}
	)
	if _, err = api.Do("/foreign/user/getInfo", nil, data); err != nil {
		return
	}
	ret = &Agent{
		ID:         data.ID,
		Phone:      data.Phone,
		Name:       data.Name,
		Sex:        data.Sex,
		TimeExpire: time.Unix(data.OverdueTime, 0),
		TimeLogin:  time.Unix(data.LastTime, 0),
	}
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

// 取群消息 gid: 群ID 群消息最少拉取30条，默认倒序返回
func (api *WkTeam) GetMsgGroup(gid string, query *Query) (ret []*MsgGroup, err error) {
	if len(gid) == 0 {
		err = contrib.ErrParamInvalid.SetVars("gid")
		return
	}
	if query == nil {
		query = &Query{}
	}
	if len(query.Account) == 0 {
		if query.Account = api.Account; len(query.Account) == 0 {
			query.Account = Settings.Account
		}
	}
	if query.Limit < 30 {
		query.Limit = 30 // 群消息最少拉取30条
	}
	if query.Params == nil {
		query.Params = make(map[string]interface{})
	}

	// 整理参数
	query.Params["account"] = gid
	if len(query.Account) > 0 {
		query.Params["my_account"] = query.Account
	}

	var (
		data []*struct {
			Uid       string `json:"wac_account"`
			Name      string `json:"wac_name"`
			NameAlias string `json:"form_name"`
			Content   string `json:"content"`
			Time      int64  `json:"create_time"`
		}
	)
	if _, err = api.Do("/foreign/group/getGroup", query, &data); err != nil {
		return
	} else {
		// 倒序结果
		for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
			data[i], data[j] = data[j], data[i]
		}
	}

	//
	for _, d := range data {
		ret = append(ret, &MsgGroup{
			Account:   api.Account,
			Gid:       gid,
			Uid:       d.Uid,
			Name:      d.Name,
			NameAlias: d.NameAlias,
			Content:   d.Content,
			Time:      time.Unix(d.Time, 0),
		})
	}
	return
}

// GetMsgUser 获取单聊信息 toUid: 好友微信唯一ID
func (api *WkTeam) GetMsgUser(toUid string, query *Query) (ret []*MsgUser, err error) {
	if len(toUid) == 0 {
		err = contrib.ErrParamUndefined.SetVars("toUid")
		return
	}
	if query == nil {
		query = &Query{}
	}
	if len(query.Account) == 0 {
		if query.Account = api.Account; len(query.Account) == 0 {
			query.Account = Settings.Account
		}
	}
	if query.Limit < 30 {
		query.Limit = 30 // 消息最少拉取30条
	}
	if query.Params == nil {
		query.Params = make(map[string]interface{})
	}
	// 整理参数
	query.Params["account"] = toUid
	if len(query.Account) > 0 {
		query.Params["my_account"] = query.Account
	}

	var (
		data []*MsgUser
	)
	if _, err = api.Do("/foreign/group/getSingle", query, &data); err != nil {
		return
	} else {
		// 倒序结果
		for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
			data[i], data[j] = data[j], data[i]
		}
		for _, d := range data {
			_ = d.init()
			d.ToUid = toUid
		}
	}
	ret = data

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

// GetGroupInfo 获取一个群信息
func (api *WkTeam) GetGroupInfo(gid string) (ret *Group, err error) {
	var (
		param = map[string]interface{}{
			"my_account": Settings.Account,
			"g_number":   gid,
		}
		data = new(Group)
	)

	if _, _err := api.Do("/foreign/group/groupInfo", &Query{Params: param}, data); _err != nil {
		err = _err
		return
	}
	ret = data

	return
}

// GetGroupInfo 获取一个用户的信息
func (api *WkTeam) GetUserInfo(account string) (ret *User, err error) {
	var (
		param = map[string]interface{}{
			"my_account": account, // 微信号
		}
		data = new(User)
	)

	if _, _err := api.Do("/foreign/message/wxInfo", &Query{Params: param}, data); _err != nil {
		err = _err
		return
	}
	ret = data
	return
}
