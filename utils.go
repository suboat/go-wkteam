package wkteam

import (
	"encoding/json"
)

//
func PubJSON(v_ interface{}) (r string) {
	r = "{}"
	if v_ == nil {
		return
	}
	switch v := v_.(type) {
	case string:
		if len(v) > 2 && v[0] == '{' && v[len(v)-1] == '}' {
			r = v
		}
		break
	case *string:
		_v := *v
		if len(_v) > 2 && _v[0] == '{' && _v[len(_v)-1] == '}' {
			r = _v
		}
		break
	default:
		// null: v_不为nil但其实是空
		if r1, _err := json.Marshal(v_); _err == nil && string(r1) != "null" {
			r = string(r1)
		}
		break
	}
	return
}
