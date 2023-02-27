package datarangers_sdk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

func Sha256Hmac(key, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func Sign(ak string, sk string, expiration int, method string, url string, params map[string]string, body string) string {
	text := canonicalRequest(method, url, params, body)
	return doSign(ak, sk, expiration, text)
}

func doSign(ak string, sk string, expiration int, text string) string {
	current := strconv.FormatInt(time.Now().Unix(), 10)
	signKeyInfo := "ak-v1/" + ak + "/" + current + "/" + strconv.Itoa(expiration)
	signKey := Sha256Hmac(sk, signKeyInfo)
	signResult := Sha256Hmac(signKey, text)
	return signKeyInfo + "/" + signResult
}

func canonicalMethod(method string) string {
	return "HTTPMethod:" + method
}

func canonicalUrl(url string) string {
	return "CanonicalURI:" + url
}

func formatKeyValue(key, value string) string {
	return key + "=" + value
}

func canonicalParam(params map[string]string) string {
	res := "CanonicalQueryString:"
	if len(params) == 0 {
		return res
	}
	for key := range params {
		res = res + formatKeyValue(key, params[key]) + "&"
	}
	return res[0 : len(res)-1]
}

func canonicalBody(body string) string {
	res := "CanonicalBody:"
	if body == "" {
		return res
	}
	return res + body
}
func canonicalRequest(method string, url string, params map[string]string, body string) string {
	cm := canonicalMethod(method)
	cu := canonicalUrl(url)
	cp := canonicalParam(params)
	cb := canonicalBody(body)
	return cm + "\n" + cu + "\n" + cp + "\n" + cb
}
