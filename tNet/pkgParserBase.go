package tNet

import "tzyNet/tINet"

func NewPkgParser[pType tINet.IPkgParser, oType tINet.IPkg](pkgBase oType) tINet.IPkgParser {
	var parser pType
	parser.SetPkgObjBase(pkgBase)
	return parser
}