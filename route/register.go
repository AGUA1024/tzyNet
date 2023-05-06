package route

import (
	"hdyx/api"
	"hdyx/api/user"
	_ "hdyx/api/user"
)

func init() {
}

func GetApiByCmd(cmd uint32) api.ApiInterface {
	highCmd := (cmd >> 16) & 0xffff
	lowCmd := cmd & 0xffff

	return masterRouteRegister[highCmd][lowCmd]
}

var masterRouteRegister = map[uint32]map[uint32]api.ApiInterface{
	0x1: userRouteRegMp,
}

var userRouteRegMp = map[uint32]api.ApiInterface{
	0x1: user.GetUserInfo_Api{},
}
