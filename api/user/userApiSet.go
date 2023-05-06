package user

import "reflect"

type GetUserInfo_Api struct {
}

func (this GetUserInfo_Api) GetFunc() any {
	return GetUserInfo
}

func (this GetUserInfo_Api) GetInType() reflect.Type {
	return reflect.TypeOf(GetUserInfo_InObj{})
}

func (this GetUserInfo_Api) GetOutType() reflect.Type {
	return reflect.TypeOf(GetUserInfo_OutObj{})
}
