package user

import "fmt"

func GetUserInfo(params *GetUserInfo_InObj) {
	fmt.Println("GetUserInfo")
	fmt.Println(params)
}
