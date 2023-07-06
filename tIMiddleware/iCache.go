package tIMiddleware

type ICache interface {
	Do(command string, args ...interface{}) (any, error)
}
