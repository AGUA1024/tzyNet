package tNet

import "tzyNet/tINet"

func NewPkgParser[pType tINet.IPkgParser]() tINet.IPkgParser {
	var p pType
	parser := p.NewParser()
	parser.SetPkgObjBase()
	return parser
}
