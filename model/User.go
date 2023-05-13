package model

import "net/http"

type User struct {
	UserId   int64       `json:"user_id"`
	NickName string      `json:"user_name"`
	RealName string      `json:"real_name"`
	AhuInfo  AhuFormInfo `json:"ahu_info"`
}

type AhuFormInfo struct {
	AhuCookie          http.Cookie `json:"ahu_cookie"`
	AhuNumber          string      `json:"ahu_number"`
	AhuPasswd          string      `json:"ahu_passwd"`
	AhuPasswdRsa       string      `json:"ahu_passwd_rsa"`
	AhuVerifyCode      string      `json:"ahu_verify_code"`
	AhuStatus          string      `json:"status"` // noLogin / logging / logged
	ET                 string      `json:"et"`
	NT                 string      `json:"nt"`
	LASTFOCUS          string      `json:"lastfocus"`
	VIEWSTATE          string      `json:"viewstate"`
	VIEWSTATEGENERATOR string      `json:"viewstategenerator"`
	EVENTTARGET        string      `json:"eventtarget"`
	EVENTARGUMENT      string      `json:"eventargument"`
}
