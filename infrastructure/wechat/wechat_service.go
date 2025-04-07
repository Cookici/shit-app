package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// WechatService 微信服务接口
type WechatService interface {
	Code2Session(code string) (*Code2SessionResponse, error)
}

// wechatService 微信服务实现
type wechatService struct {
	appID     string
	appSecret string
}

// NewWechatService 创建微信服务
func NewWechatService(appID, appSecret string) WechatService {
	return &wechatService{
		appID:     appID,
		appSecret: appSecret,
	}
}

// Code2SessionResponse 微信登录返回结果
type Code2SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// Code2Session 微信登录，用code换取openid和session_key
func (s *wechatService) Code2Session(code string) (*Code2SessionResponse, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		s.appID, s.appSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Code2SessionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, errors.New(result.ErrMsg)
	}

	return &result, nil
}
