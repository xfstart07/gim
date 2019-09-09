// Author: xufei
// Date: 2019-09-09 09:40

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gim/internal/lg"
	"gim/model"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

var (
	errP2PSendFail = errors.New("私聊消息发送失败")
)

func (c *Client) sendP2PMsg(msgReq model.P2PReq) error {
	url := fmt.Sprintf("http://%s:%s/sendP2PMsg", GetConfig().ServerIP, GetConfig().ServerPort)

	msgBody, _ := json.Marshal(msgReq)
	resp, err := http.Post(url, ContextTypeJSON, bytes.NewBuffer(msgBody))
	if err != nil {
		lg.Logger().Error("私聊消息发送失败", zap.Error(err))
		return errP2PSendFail
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	lg.Logger().Info("发送结果" + string(respBody))

	return nil
}
