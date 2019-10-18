package wkteam

import (
	"time"
)

// 定义一些常用的快捷方法

// GetMsgGroupSince 获取某时刻及以后的所有群聊消息, 如循环获取数据过程中出错, ret也会返回目前为止收集到的数据
func (api *WkTeam) GetMsgGroupSince(gid string, gte time.Time) (ret []*MsgGroup, err error) {
	var (
		limit = 30
		query = &Query{
			Limit: limit,
		}
		since = time.Now()
		total = 0
	)

	// 循环取数据
	for since.After(gte) {
		var data []*MsgGroup
		if data, err = api.GetMsgGroup(gid, query); err != nil {
			return
		} else if len(data) == 0 {
			// 没有了
			break
		}
		total += len(data)
		for _, d := range data {
			if d.Time.Before(gte) == false {
				ret = append(ret, d)
			}
		}

		// break: 最后一页了
		if len(data) < query.Limit || query.Total == total {
			break
		}

		// next
		if len(data) > 0 {
			since = data[len(data)-1].Time
		}
		query.Skip += query.Limit
	}

	// debug
	api.Log.Debugf(`[api-shortcut] GetMsgGroupSince get %d/%d`, len(ret), total)
	return
}

// GetMsgUserSince 获取某时刻及以后的所有单聊消息, 如循环获取数据过程中出错, ret也会返回目前为止收集到的数据
func (api *WkTeam) GetMsgUserSince(targetUid string, gte time.Time) (ret []*MsgUser, err error) {
	var (
		limit = 30
		query = &Query{
			Limit: limit,
		}
		since = time.Now()
		total = 0
	)

	// 循环取数据
	for since.After(gte) {
		var data []*MsgUser
		if data, err = api.GetMsgUser(targetUid, query); err != nil {
			return
		} else if len(data) == 0 {
			// 没有了
			break
		}
		total += len(data)
		for _, d := range data {
			if d.Time.Before(gte) == false {
				ret = append(ret, d)
			}
		}

		// break: 最后一页了
		if len(data) < query.Limit || query.Total == total {
			break
		}

		// next
		if len(data) > 0 {
			since = data[len(data)-1].Time
		}
		query.Skip += query.Limit
	}

	// debug
	api.Log.Debugf(`[api-shortcut] GetMsgUserSince get %d/%d`, len(ret), total)
	return
}
