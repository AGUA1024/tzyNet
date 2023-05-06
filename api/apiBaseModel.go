package api

import "reflect"

type ApiInterface interface {
	GetFunc() any
	GetInType() reflect.Type
	GetOutType() reflect.Type
}
