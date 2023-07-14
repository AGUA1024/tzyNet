package tNet

import "tzyNet/tINet"

func NewPkgParser[pType tINet.IMsgParser]() tINet.IMsgParser {
	var p pType
	parser := p.NewParser()
	parser.SetPkgObjBase()
	return parser
}
