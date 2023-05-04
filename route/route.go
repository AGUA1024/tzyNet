package route

func anyRoute() {

}

var RouteRegister = map[int32]func(){
	1: anyRoute,
}
