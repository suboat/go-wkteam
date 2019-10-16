package wkteam

import (
	"github.com/stretchr/testify/require"

	"testing"
)

//
func Test_Agent(t *testing.T) {
	as := require.New(t)
	d := new(Agent)
	as.Nil(d.Parse(`{"uid":1111,"phone":"123456","name":"有限公司","sex":"未知","email":"","overdue_time":1573488000,"last_time":1571119420}`))
	t.Log(d)
}
