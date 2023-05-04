/**
 *	Copyright 2020 Beijing Volcano Engine Technology Co., Ltd.
 *	Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *  http://www.apache.org/licenses/LICENSE-2.0
 *	Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

package datarangers_sdk

type ServerSdkEventMessage struct {
	ApiRequestId  *string    `json:"_api_request_id,omitempty"`
	HttpHost      *string    `json:"_http_host,omitempty"`
	ClientIp      *string    `json:"client_ip,omitempty"`
	DeviceId      *int64     `json:"device_id,omitempty"`
	Header        *Header    `json:"header"`
	TimeSync      *TimeSync  `json:"time_sync,omitempty"`
	TraceId       *string    `json:"trace_id,omitempty"`
	UserId        *string    `json:"user_id,omitempty"`
	EventV3       []*EventV3 `json:"event_v3,omitempty"`
	ApiTime       *int64     `json:"api_time,omitempty"`
	HttpUserAgent *string    `json:"http_user_agent,omitempty"`
	UserIsAuth    *int64     `json:"user_is_auth,omitempty"`
	Ut            *int64     `json:"ut,omitempty"`
	AppType       *string    `json:"app_type"`
	FormatName    *string    `json:"_format_name,omitempty"`
	AppId         *int64     `json:"app_id,omitempty"`
	UserUniqueId  *string    `json:"user_unique_id,omitempty"`
	MessageType   string     `json:"-"`
}

type Header struct {
	Aid                *int64                 `json:"aid,omitempty"`
	AppLanguage        *string                `json:"app_language,omitempty"`
	AppName            *string                `json:"app_name,omitempty"`
	AppRegion          *string                `json:"app_region,omitempty"`
	AppVersion         *string                `json:"app_version,omitempty"`
	AppVersionMinor    *string                `json:"app_version_minor,omitempty"`
	BuildSerial        *string                `json:"build_serial,omitempty"`
	Carrier            *string                `json:"carrier,omitempty"`
	Channel            *string                `json:"channel,omitempty"`
	Clientudid         *string                `json:"clientudid,omitempty"`
	CpuAbi             *string                `json:"cpu_abi,omitempty"`
	Custom             map[string]interface{} `json:"custom,omitempty"`
	DeviceId           *int64                 `json:"device_id,omitempty"`
	DeviceBrand        *string                `json:"device_brand,omitempty"`
	DeviceManufacturer *string                `json:"device_manufacturer,omitempty"`
	DeviceModel        *string                `json:"device_model,omitempty"`
	DeviceType         *string                `json:"device_type,omitempty"`
	DisplayName        *string                `json:"display_name,omitempty"`
	DisplayDensity     *string                `json:"display_density,omitempty"`
	DensityDpi         *string                `json:"density_dpi,omitempty"`
	Idfa               *string                `json:"idfa,omitempty"`
	InstallId          *uint64                `json:"install_id,omitempty"`
	Language           *string                `json:"language,omitempty"`
	Openudid           *string                `json:"openudid,omitempty"`
	VendorId           *string                `json:"vendor_id,omitempty"`
	Os                 *string                `json:"os,omitempty"`
	OsVersion          *string                `json:"os_version,omitempty"`
	OsApi              *string                `json:"os_api,omitempty"`
	Package            *string                `json:"package,omitempty"`
	Region             *string                `json:"region,omitempty"`
	SdkVersion         *string                `json:"sdk_version,omitempty"`
	Timezone           *int32                 `json:"timezone,omitempty"`
	TzOffset           *int64                 `json:"tz_offset,omitempty"`
	TzName             *string                `json:"tz_name,omitempty"`
	Udid               *string                `json:"udid,omitempty"`
	UserUniqueId       *string                `json:"user_unique_id,omitempty"`
}

type TimeSync struct {
	LocalTime  *int64 `json:"local_time,omitempty"`
	ServerTime *int64 `json:"server_time,omitempty"`
}

type EventV3 struct {
	Event         string                 `json:"event,omitempty"`
	Params        map[string]interface{} `json:"params,omitempty"`
	LocalTimeMs   *int64                 `json:"local_time_ms,omitempty"`
	Datetime      *string                `json:"datetime,omitempty"`
	EventId       *string                `json:"event_id,omitempty"`
	AbSdkVersion  *string                `json:"ab_sdk_version,omitempty"`
	SessionId     *string                `json:"session_id,omitempty"`
	UserId        *string                `json:"user_id,omitempty"`
	TeaEventIndex *int64                 `json:"tea_event_index,omitempty"`
}

type Item struct {
	ItemName *string
	ItemId   *string
}

type SaasServerEvent struct {
	Event         string  `json:"event,omitempty"`
	Params        *string `json:"params,omitempty"`
	SessionId     *string `json:"session_id,omitempty"`
	LocalTimeMs   *int64  `json:"local_time_ms,omitempty"`
	Datetime      *string `json:"datetime,omitempty"`
	TeaEventIndex *int64  `json:"tea_event_index,omitempty"`
	AbSdkVersion  *string `json:"ab_sdk_version,omitempty"`
}

type User struct {
	UserUniqueId *string `json:"user_unique_id,omitempty"`
}

type DefaultSaasServerAppMessage struct {
	User   *User              `json:"user,omitempty"`
	Header *Header            `json:"header,omitempty"`
	Events []*SaasServerEvent `json:"events,omitempty"`
}

type SaasServerAppMessage struct {
	*DefaultSaasServerAppMessage
}

type Attribute struct {
	Name      *string      `json:"name,omitempty"`
	Value     *interface{} `json:"value,omitempty"`
	Operation string       `json:"operation,omitempty"`
}

type SaasProfileAppMessage struct {
	Attributes []*Attribute `json:"attributes,omitempty"`
}

type SaasItemAppMessage struct {
	Attributes []*Attribute `json:"attributes,omitempty"`
}

func (spam *SaasProfileAppMessage) AddAttribute(name string, value interface{}, method string) {
	operation, ok := profileOperationMap[method]
	if !ok {
		panic("Not support method: " + method)
	}

	attribute := &Attribute{
		Name:      &name,
		Value:     &value,
		Operation: operation,
	}
	spam.Attributes = append(spam.Attributes, attribute)
}

func (siam *SaasItemAppMessage) AddAttribute(name string, value interface{}, method string) {
	operation, ok := itemOperationMap[method]
	if !ok {
		panic("Not support method: " + method)
	}

	attribute := &Attribute{
		Name:      &name,
		Value:     &value,
		Operation: operation,
	}
	siam.Attributes = append(siam.Attributes, attribute)
}
