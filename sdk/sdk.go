package sdk

import (
	"bytes"
	"io"
	"net/http"
)

type SdkInterface interface {
	GetUrl() string
	GetSign(strArgs ...string) string
}

func SdkRequest(sdk SdkInterface, byteReqBuf []byte) []byte {
	req, err := http.NewRequest(http.MethodPost, sdk.GetUrl(), bytes.NewReader(byteReqBuf))
	if err != nil {
		//common.Logger.GameErrorLog(ctx, common.ERR_PLAYERINFO_REQUEST, "建立http会话失败")
		return nil
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//common.Logger.GameErrorLog(ctx, common.ERR_PLAYERINFO_REQUEST, "向不鸽平台请求玩家信息失败")
		return nil
	}
	defer resp.Body.Close()

	// 解析响应
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return body
}
