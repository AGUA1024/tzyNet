package model

import (
	"reflect"
	"tzyNet/tCommon"
)

const newFuncName = "NewActModel"

var actModelRegister = map[uint32]ActModelInterface{
	1: &Act1Model{},
}

func NewActModel(ctx *tCommon.ConContext, actId uint32) ActModelInterface {
	model, ok := actModelRegister[actId]
	if !ok {
		return nil
	}

	// 获取结构体类型的反射值
	value := reflect.ValueOf(model)
	newActFunc := value.MethodByName(newFuncName)

	argsValues := []reflect.Value{
		reflect.ValueOf(ctx),
	}
	// 调用方法
	results := newActFunc.Call(argsValues)

	// 处理返回值
	var ret ActModelInterface
	if len(results) > 0 {
		retValue := results[0]
		iface := retValue.Interface()
		if ret, ok = iface.(ActModelInterface); !ok {
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
