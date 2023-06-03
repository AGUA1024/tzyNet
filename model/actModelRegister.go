package model

import (
	"hdyx/common"
	"reflect"
)

const newFuncName = "NewActModel"

var actModelRegister = map[uint32]ActBaseInterface{
	1: &Act1Model{},
}

func NewActModel(ctx *common.ConContext, actId uint32) ActBaseInterface {
	model, ok := actModelRegister[actId]
	if !ok {
		return nil
	}

	// 获取结构体类型的反射值
	modelType := reflect.TypeOf(model)
	modelPtr := reflect.New(modelType)
	elem := modelPtr.Elem()

	// 获取方法的反射值
	method := elem.MethodByName(newFuncName)

	argsValues := []reflect.Value{
		reflect.ValueOf(ctx),
	}
	// 调用方法
	results := method.Call(argsValues)

	// 处理返回值
	var ret ActBaseInterface
	if len(results) > 0 {
		retValue := results[0]
		iface := retValue.Interface()
		if ret, ok = iface.(ActBaseInterface); !ok {
			return nil
		}
	}

	return ret
}

func GetActCfg(actId uint32) *GameCfg {
	model, ok := actModelRegister[actId]
	if !ok {
		return nil
	}
	return model.GetActCfg()
}
