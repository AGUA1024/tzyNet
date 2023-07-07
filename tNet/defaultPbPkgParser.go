package tNet

import (
	"errors"
	"google.golang.org/protobuf/proto"
	"tzyNet/tINet"
	"tzyNet/tNet/ioBuf"
)

type DefaultPbPkgParser struct {
	pkgBaseObj tINet.IPkg
}

type DefaultPbPkg struct {
	dataSrc *ioBuf.ClientBuf
}

func (this *DefaultPbPkgParser) SetPkgObjBase(pkgBase tINet.IPkg) {
	this.pkgBaseObj = pkgBase
}

func (this *DefaultPbPkgParser) Marshal(obj any) ([]byte, error) {
	pbObj, ok := obj.(proto.Message)
	if !ok {
		return nil, errors.New("Marsahl_Error_invalid_ojb_type")
	}

	byteMsg, err := proto.Marshal(pbObj)
	return byteMsg, err
}

func (this *DefaultPbPkgParser) UnMarshal(byteMsg []byte) (tINet.IPkg, error) {
	base := this.pkgBaseObj
	pbBase := base.(proto.Message)
	err := proto.Unmarshal(byteMsg, pbBase)
	if err != nil {
		return nil, err
	}

	return pbBase, nil
}

func (this *DefaultPbPkg) GetRouteCmd() uint32 {
	return this.dataSrc.CmdMerge
}

func (this *DefaultPbPkg) GetData() []byte {
	return this.dataSrc.Data
}
