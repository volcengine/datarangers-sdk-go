package event_collect

type dancemsg struct {
	Api_request_id  *string   `json:"_api_request_id,omitempty"`
	Http_host       *string   `json:"_http_host,omitempty"`
	Client_ip       *string   `json:"client_ip,omitempty"`
	Device_id       *int64    `json:"device_id,omitempty"`
	Header          *header   `json:"header"`
	TimeSync        *timeSync `json:"time_sync,omitempty"`
	Trace_id        *string   `json:"trace_id,omitempty"`
	U_id            *string   `json:"u_id,omitempty"`
	User_id         *string   `json:"user_id,omitempty"`
	Event_v3        []*items  `json:"event_v3,omitempty"`
	Api_time        *int64    `json:"api_time,omitempty"`
	Http_user_agent *string   `json:"http_user_agent,omitempty"`
	User_is_auth    *int64    `json:"user_is_auth,omitempty"`
	Ut              *int64    `json:"ut,omitempty"`
	App_type        *string   `json:"app_type"`
	Format_name     *string   `json:"_format_name,omitempty"`
	App_id          *int64    `json:"app_id"`
	User_unique_id  *string   `json:"user_unique_id,omitempty"`
}

type header struct {
	Aid                 *int64                 `json:"aid"`
	App_language        *string                `json:"app_language,omitempty"`
	App_name            *string                `json:"app_name,omitempty"`
	App_region          *string                `json:"app_region,omitempty"`
	App_version         *string                `json:"app_version,omitempty"`
	App_version_minor   *string                `json:"app_version_minor,omitempty"`
	Appkey              *string                `json:"appkey,omitempty"`
	Build_serial        *string                `json:"build_serial,omitempty"`
	Carrier             *string                `json:"carrier,omitempty"`
	Channel             *string                `json:"channel,omitempty"`
	Clientudid          *string                `json:"clientudid,omitempty"`
	Cpu_abi             *string                `json:"cpu_abi,omitempty"`
	Custom              map[string]interface{} `json:"custom,omitempty"`
	Device_id           *int64                 `json:"device_id,omitempty"`
	Device_brand        *string                `json:"device_brand,omitempty"`
	Device_manufacturer *string                `json:"device_manufacturer,omitempty"`
	Device_model        *string                `json:"device_model,omitempty"`
	Device_type         *string                `json:"device_type,omitempty"`
	Display_name        *string                `json:"display_name,omitempty"`
	Display_density     *string                `json:"display_density,omitempty"`
	Density_dpi         *string                `json:"density_dpi,omitempty"`
	Idfa                *string                `json:"idfa,omitempty"`
	Install_id          *uint64                `json:"install_id,omitempty"`
	Language            *string                `json:"language,omitempty"`
	Openudid            *string                `json:"openudid,omitempty"`
	Os                  *string                `json:"os,omitempty"`
	OsVersion           *string                `json:"os_version,omitempty"`
	OsApi               *string                `json:"os_api,omitempty"`
	Package             *string                `json:"package,omitempty"`
	Region              *string                `json:"region,omitempty"`
	Sdk_version         *string                `json:"sdk_version,omitempty"`
	Timezone            *int64                 `json:"timezone,omitempty"`
	TzOffset            *int64                 `json:"tz_offset,omitempty"`
	TzName              *string                `json:"tz_name,omitempty"`
	Udid                *string                `json:"udid,omitempty"`
	User_unique_id      *string                `json:"user_unique_id,omitempty"`
}

type timeSync struct {
	Local_time  *int64 `json:"local_time,omitempty"`
	Server_time *int64 `json:"server_time,omitempty"`
}

type items struct {
	Datetime        *string     `json:"datetime,omitempty"`
	Event           *string     `json:"event,omitempty"`
	EventId         *string     `json:"event_id,omitempty"`
	LocalTimeMs     *int64      `json:"local_time_ms,omitempty"`
	SessionId       *string     ` json:"session_id,omitempty"`
	UserId          *string     `json:"user_id,omitempty"`
	Tea_event_index *int64      `json:"tea_event_index,omitempty"`
	Params          interface{} `json:"params,omitempty"`
}
