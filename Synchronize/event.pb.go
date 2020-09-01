package Synchronize

//import proto "github.com/golang/protobuf/proto"
import (
	"github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type user struct {
	UserUniqueId     *string `protobuf:"bytes,1,req,name=user_unique_id" json:"user_unique_id,omitempty"`
	UserType         *uint32 `protobuf:"varint,2,opt,name=user_type" json:"user_type,omitempty"`
	UserId           *uint64 `protobuf:"varint,3,opt,name=user_id" json:"user_id,omitempty"`
	UserIsAuth       *bool   `protobuf:"varint,4,opt,name=user_is_auth" json:"user_is_auth,omitempty"`
	UserIsLogin      *bool   `protobuf:"varint,5,opt,name=user_is_login" json:"user_is_login,omitempty"`
	DeviceId         *uint64 `protobuf:"varint,6,opt,name=device_id" json:"device_id,omitempty"`
	WebId            *uint64 `protobuf:"varint,7,opt,name=web_id" json:"web_id,omitempty"`
	IpAddrId         *uint64 `protobuf:"varint,8,opt,name=ip_addr_id" json:"ip_addr_id,omitempty"`
	Ssid             *string `protobuf:"bytes,9,opt,name=ssid" json:"ssid,omitempty"`
	OpenUdid         *string `protobuf:"bytes,10,opt,name=open_udid" json:"open_udid,omitempty"`
	Udid             *string `protobuf:"bytes,11,opt,name=udid" json:"udid,omitempty"`
	Idfa             *string `protobuf:"bytes,12,opt,name=idfa" json:"idfa,omitempty"`
	Idfv             *string `protobuf:"bytes,13,opt,name=idfv" json:"idfv,omitempty"`
	BuildSerial      *string `protobuf:"bytes,14,opt,name=build_serial" json:"build_serial,omitempty"`
	ClientUdid       *string `protobuf:"bytes,15,opt,name=client_udid" json:"client_udid,omitempty"`
	Mc               *string `protobuf:"bytes,16,opt,name=mc" json:"mc,omitempty"`
	SerialNumber     *string `protobuf:"bytes,17,opt,name=serial_number" json:"serial_number,omitempty"`
	IsUpgradeUser    *bool   `protobuf:"varint,18,opt,name=is_upgrade_user" json:"is_upgrade_user,omitempty"`
	Oaid             *string `protobuf:"bytes,19,opt,name=oaid" json:"oaid,omitempty"`
	Cdid             *string `protobuf:"bytes,20,opt,name=cdid" json:"cdid,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

type header struct {
	Headers           *string                `protobuf:"bytes,1,opt,name=headers" json:"headers,omitempty"`
	AppId             *uint32                `protobuf:"varint,2,req,name=app_id" json:"app_id,omitempty"`
	AppName           *string                `protobuf:"bytes,3,opt,name=app_name" json:"app_name,omitempty"`
	AppInstallId      *uint64                `protobuf:"varint,4,opt,name=app_install_id" json:"app_install_id,omitempty"`
	AppPackage        *string                `protobuf:"bytes,5,opt,name=app_package" json:"app_package,omitempty"`
	AppChannel        *string                `protobuf:"bytes,6,opt,name=app_channel" json:"app_channel,omitempty"`
	AppVersion        *string                `protobuf:"bytes,7,opt,name=app_version" json:"app_version,omitempty"`
	OsName            *string                `protobuf:"bytes,8,opt,name=os_name" json:"os_name,omitempty"`
	OsVersion         *string                `protobuf:"bytes,9,opt,name=os_version" json:"os_version,omitempty"`
	DeviceModel       *string                `protobuf:"bytes,10,opt,name=device_model" json:"device_model,omitempty"`
	AbClient          *string                `protobuf:"bytes,11,opt,name=ab_client" json:"ab_client,omitempty"`
	AbVersion         *string                `protobuf:"bytes,12,opt,name=ab_version" json:"ab_version,omitempty"`
	TrafficType       *string                `protobuf:"bytes,13,opt,name=traffic_type" json:"traffic_type,omitempty"`
	UtmSource         *string                `protobuf:"bytes,14,opt,name=utm_source" json:"utm_source,omitempty"`
	UtmMedium         *string                `protobuf:"bytes,15,opt,name=utm_medium" json:"utm_medium,omitempty"`
	UtmCampaign       *string                `protobuf:"bytes,16,opt,name=utm_campaign" json:"utm_campaign,omitempty"`
	ClientIp          *string                `protobuf:"bytes,17,opt,name=client_ip" json:"client_ip,omitempty"`
	DeviceBrand       *string                `protobuf:"bytes,18,opt,name=device_brand" json:"device_brand,omitempty"`
	OsApi             *uint32                `protobuf:"varint,19,opt,name=os_api" json:"os_api,omitempty"`
	Access            *string                `protobuf:"bytes,20,opt,name=access" json:"access,omitempty"`
	Language          *string                `protobuf:"bytes,21,opt,name=language" json:"language,omitempty"`
	Region            *string                `protobuf:"bytes,22,opt,name=region" json:"region,omitempty"`
	AppLanguage       *string                `protobuf:"bytes,23,opt,name=app_language" json:"app_language,omitempty"`
	AppRegion         *string                `protobuf:"bytes,24,opt,name=app_region" json:"app_region,omitempty"`
	CreativeId        *uint64                `protobuf:"varint,25,opt,name=creative_id" json:"creative_id,omitempty"`
	AdId              *uint64                `protobuf:"varint,26,opt,name=ad_id" json:"ad_id,omitempty"`
	CampaignId        *uint64                `protobuf:"varint,27,opt,name=campaign_id" json:"campaign_id,omitempty"`
	LogType           *string                `protobuf:"bytes,28,opt,name=log_type" json:"log_type,omitempty"`
	Rnd               *string                `protobuf:"bytes,29,opt,name=rnd" json:"rnd,omitempty"`
	Platform          *string                `protobuf:"bytes,30,opt,name=platform" json:"platform,omitempty"`
	SdkVersion        *string                `protobuf:"bytes,31,opt,name=sdk_version" json:"sdk_version,omitempty"`
	Province          *string                `protobuf:"bytes,32,opt,name=province" json:"province,omitempty"`
	City              *string                `protobuf:"bytes,33,opt,name=city" json:"city,omitempty"`
	Timezone          *int32                 `protobuf:"varint,34,opt,name=timezone" json:"timezone,omitempty"`
	TzOffset          *int32                 `protobuf:"varint,35,opt,name=tz_offset" json:"tz_offset,omitempty"`
	TzName            *string                `protobuf:"bytes,36,opt,name=tz_name" json:"tz_name,omitempty"`
	SimRegion         *string                `protobuf:"bytes,37,opt,name=sim_region" json:"sim_region,omitempty"`
	Carrier           *string                `protobuf:"bytes,38,opt,name=carrier" json:"carrier,omitempty"`
	Resolution        *string                `protobuf:"bytes,39,opt,name=resolution" json:"resolution,omitempty"`
	Browser           *string                `protobuf:"bytes,50,opt,name=browser" json:"browser,omitempty"`
	BrowserVersion    *string                `protobuf:"bytes,51,opt,name=browser_version" json:"browser_version,omitempty"`
	Referrer          *string                `protobuf:"bytes,52,opt,name=referrer" json:"referrer,omitempty"`
	ReferrerHost      *string                `protobuf:"bytes,53,opt,name=referrer_host" json:"referrer_host,omitempty"`
	ScreenHeight      *int32                 `protobuf:"varint,54,opt,name=screen_height" json:"screen_height,omitempty"`
	ScreenWidth       *int32                 `protobuf:"varint,55,opt,name=screen_width" json:"screen_width,omitempty"`
	Tz                *float32               `protobuf:"fixed32,56,opt,name=tz" json:"tz,omitempty"`
	AppVersionMinor   *string                `protobuf:"bytes,57,opt,name=app_version_minor" json:"app_version_minor,omitempty"`
	CarrierRegion     *string                `protobuf:"bytes,58,opt,name=carrier_region" json:"carrier_region,omitempty"`
	ProductName       *string                `protobuf:"bytes,59,opt,name=product_name" json:"product_name,omitempty"`
	ProductId         *uint32                `protobuf:"varint,60,opt,name=product_id" json:"product_id,omitempty"`
	Custom            map[string]interface{} `protobuf:"bytes,61,opt,name=custom" json:"custom,omitempty"`
	UpdateVersionCode *uint32                `protobuf:"varint,62,opt,name=update_version_code" json:"update_version_code,omitempty"`
	AbSdkVersion      *string                `protobuf:"bytes,63,opt,name=ab_sdk_version" json:"ab_sdk_version,omitempty"`
	UserAgent         *string                `protobuf:"bytes,64,opt,name=user_agent" json:"user_agent,omitempty"`
	ClientPort        *uint32                `protobuf:"varint,65,opt,name=client_port" json:"client_port,omitempty"`
	DataCenter        *string                `protobuf:"bytes,66,opt,name=data_center" json:"data_center,omitempty"`
	OriginAppId       *uint32                `protobuf:"varint,67,opt,name=origin_app_id" json:"origin_app_id,omitempty"`
	OriginAppName     *string                `protobuf:"bytes,68,opt,name=origin_app_name" json:"origin_app_name,omitempty"`
	UtmTerm           *string                `protobuf:"bytes,69,opt,name=utm_term" json:"utm_term,omitempty"`
	UtmContent        *string                `protobuf:"bytes,70,opt,name=utm_content" json:"utm_content,omitempty"`
	XXX_unrecognized  []byte                 `json:"-"`

	AppAppId       *uint32 `protobuf:"varint,2,req,name=app_id" json:"aid,omitempty"`
	DeviceId       *uint64 `protobuf:"varint,26,opt,name=ad_id" json:"device_id,omitempty"`
	SsId           *string `protobuf:"varint,26,opt,name=ad_id" json:"ssid,omitempty"`
	User_unique_id *string `protobuf:"varint,26,opt,name=ad_id" json:"user_unique_id,omitempty"`
	Install_id     *uint64 `protobuf:"varint,26,opt,name=ad_id" json:"install_id,omitempty"`
	Web_id         *uint64 `protobuf:"varint,26,opt,name=ad_id" json:"web_id,omitempty"`
}

type event struct {
	Event            *string     `protobuf:"bytes,3,req,name=event" json:"event,omitempty"`
	Time             *uint32     `protobuf:"varint,4,req,name=time" json:"time,omitempty"`
	Params           interface{} `protobuf:"bytes,5,req,name=params" json:"params,omitempty"`
	SessionId        *string     `protobuf:"bytes,6,opt,name=session_id" json:"session_id,omitempty"`
	LocalTimeMs      *uint64     `protobuf:"varint,7,opt,name=local_time_ms" json:"local_time_ms,omitempty"`
	XStagingFlag     *bool       `protobuf:"varint,8,opt,name=_staging_flag" json:"_staging_flag,omitempty"`
	SeqId            *uint32     `protobuf:"varint,9,opt,name=seq_id" json:"seq_id,omitempty"`
	EventId          *uint64     `protobuf:"varint,10,opt,name=event_id" json:"event_id,omitempty"`
	XXX_unrecognized []byte      `json:"-"`

	Datetime     *string `protobuf:"bytes,3,req,name=event" json:"datetime,omitempty"`
	UserId       *string `protobuf:"bytes,3,req,name=event" json:"user_id,omitempty"`
	Localtime_ms *uint64 `protobuf:"bytes,3,req,name=event" json:"localtime_ms,omitempty"`
}

type marioEvents struct {
	Caller           *string  `protobuf:"bytes,1,req,name=caller" json:"caller,omitempty"`
	ServerTime       *uint32  `protobuf:"varint,2,req,name=server_time" json:"server_time,omitempty"`
	User             *user    `protobuf:"bytes,3,req,name=user" json:"user,omitempty"`
	Header           *header  `protobuf:"bytes,4,req,name=header" json:"header,omitempty"`
	Events           []*event `protobuf:"bytes,5,rep,name=events" json:"events,omitempty"`
	TraceId          *string  `protobuf:"bytes,6,opt,name=trace_id" json:"trace_id,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
	AppEvents        []*event `protobuf:"bytes,5,rep,name=events" json:"event_v3,omitempty"`
	//Launchs          []*Launch    ` json:"launch,omitempty"`
	//Terminates       []*Terminate `json:"terminate,omitempty"`
}

//type Launch struct {
//	Datetime     *string `json:"datetime,omitempty"`
//	Localtime_ms *uint64 `json:"local_time_ms,omitempty"`
//	SessionId    *string `json:"session_id,omitempty"`
//	TeaEvent     *uint64 `json:"tea_event_index,omitempty"`
//}
//
//type Terminate struct {
//	Datetime     *string `json:"datetime,omitempty"`
//	Localtime_ms *uint64 `json:"local_time_ms,omitempty"`
//	SessionId    *string `json:"session_id,omitempty"`
//	TeaEvent     *uint64 `json:"tea_event_index,omitempty"`
//	Duration     *uint64 `json:"duration,omitempty"`
//}
