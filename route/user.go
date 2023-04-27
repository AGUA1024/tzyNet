package route

import (
	"hdyx/api/user"
)

func init() {
	videoGroup := R.Group("/ws")
	{
		videoGroup.GET("/ws", user.Ping)
	}
}
