package tNet

import (
	"errors"
	"google.golang.org/protobuf/proto"
	"tzyNet/tINet"
	"tzyNet/tNet/ioBuf"
)

type DefaultPbPkgParser struct {
	PkgBaseObj tINet.IMsg
}

type DefaultPbPkg struct {
	DataSrc *ioBuf.ClientBuf
}

func (this *DefaultPbPkg) GetServName() string {
	//TODO implement me
	panic("implement me")
}

func (this *DefaultPbPkg) Serialize() []byte {
	//TODO implement me
	panic("implement me")
}

func (this *DefaultPbPkg) Deserialize() tINet.IMsg {
	//TODO implement me
	panic("implement me")
}

func (this *DefaultPbPkg) SetDataSrc(src any) {
	iobuf, _ := src.(ioBuf.ClientBuf)
	this.DataSrc = &iobuf
}

func (this *DefaultPbPkgParser) NewParser() tINet.IMsgParser {
	return &DefaultPbPkgParser{PkgBaseObj: nil}
}

func (this *DefaultPbPkgParser) SetPkgObjBase() {
	this.PkgBaseObj = &DefaultPbPkg{DataSrc: &ioBuf.ClientBuf{}}
}

func (this *DefaultPbPkgParser) Marshal(obj any) ([]byte, error) {
	pbObj, ok := obj.(proto.Message)
	if !ok {
		return nil, errors.New("Marsahl_Error_invalid_ojb_type")
	}

	return proto.Marshal(pbObj)
}

func (this *DefaultPbPkgParser) UnMarshal(byteMsg []byte) (tINet.IMsg, error) {
	base := this.PkgBaseObj

	pb := ioBuf.ClientBuf{}
	err := proto.Unmarshal(byteMsg, &pb)
	if err != nil {
		return nil, err
	}

	base.SetDataSrc(pb)
	return base, nil
}

func (this *DefaultPbPkg) GetRouteCmd() uint32 {
	return this.DataSrc.CmdMerge
}

func (this *DefaultPbPkg) GetData() []byte {
	return this.DataSrc.Data
}
