package asynsdk

type message struct {
	isapp      bool                   `json:"isapp,omitempty"`
	appid      uint32                 `json:"appid,omitempty"`
	uuid       string                 `json:"uuid,omitempty"`
	eventname  string                 `json:"eventname,omitempty"`
	eventParam map[string]interface{} `json:"eventParam,omitempty"`
}
