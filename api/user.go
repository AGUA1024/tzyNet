package api

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"hdyx/api/protobuf"
	"hdyx/common"
	ioBuf2 "hdyx/net/ioBuf"
)

func GetUserInfo(params []byte) {
	parmObj := getParamObj[*protobuf.GetUserInfo_InObj](params, &protobuf.GetUserInfo_InObj{})

	fmt.Println(parmObj)
	fmt.Println(parmObj.Uid)

	//modelObj := model.NewAct1Model(3)
	//
	//act1Model := modelObj.(*model.Act1Model)

	//fmt.Println(act1Model)
	//fmt.Println("GetUserInfo")
	//fmt.Println(params)

	outPut := &ioBuf2.OutPutBuf{
		CmdCode:        12,
		ProtocolSwitch: 22,
		CmdMerge:       35,
		ResponseStatus: 67,
		ValidMsg:       "hello world",
		Data:           nil,
	}

	msg, _ := proto.Marshal(outPut)

	// 写入ws数据
	common.GameMaster.OutMsg(1, msg)

	common.Logger.ErrorLog("HAHAH")

}
