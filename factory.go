package datarangers_sdk

import "encoding/json"

func CreateDefaultSaasServerAppMessage(ssem *ServerSdkEventMessage) *DefaultSaasServerAppMessage {
	user := &User{
		UserUniqueId: ssem.UserUniqueId,
	}
	dssam := &DefaultSaasServerAppMessage{
		User:   user,
		Header: ssem.Header,
	}
	if ssem.EventV3 != nil {
		for _, eventV3 := range ssem.EventV3 {
			sse := createSaasServerEvent(eventV3)
			dssam.Events = append(dssam.Events, sse)
		}
	}
	return dssam
}

func CreateSaasServerAppMessage(ssem *ServerSdkEventMessage) *SaasServerAppMessage {
	dssam := CreateDefaultSaasServerAppMessage(ssem)
	ssam := &SaasServerAppMessage{
		DefaultSaasServerAppMessage: dssam,
	}
	ssam.Header.Aid = nil
	return ssam
}

func CreateSaasProfileAppMessage(ssem *ServerSdkEventMessage) *SaasProfileAppMessage {
	spam := &SaasProfileAppMessage{}
	if ssem.EventV3 != nil {
		for _, eventV3 := range ssem.EventV3 {
			mParams := eventV3.Params
			for k, v := range mParams {
				if k == "item_id" || k == "item_name" {
					continue
				}
				spam.AddAttribute(k, v, eventV3.Event)
			}
		}
	}
	return spam
}

func CreateSaasItemAppMessage(eventV3 *EventV3) *SaasItemAppMessage {
	siam := &SaasItemAppMessage{}
	mParams := eventV3.Params
	for k, v := range mParams {
		if k == "item_id" || k == "item_name" {
			continue
		}
		siam.AddAttribute(k, v, eventV3.Event)
	}
	return siam
}

func createSaasServerEvent(eventV3 *EventV3) *SaasServerEvent {
	params, _ := json.Marshal(eventV3.Params)
	sParams := string(params)
	sse := &SaasServerEvent{
		Event:         eventV3.Event,
		SessionId:     eventV3.SessionId,
		LocalTimeMs:   eventV3.LocalTimeMs,
		Datetime:      eventV3.Datetime,
		AbSdkVersion:  eventV3.AbSdkVersion,
		TeaEventIndex: eventV3.TeaEventIndex,
		Params:        &sParams,
	}
	return sse
}

func CreateSaasNativeAppMessage(ssem ServerSdkEventMessage) *ServerSdkEventMessage {
	snam := ssem
	snam.AppId = nil
	snam.Header.Aid = nil
	return &snam
}

func CreateSaasNativeAppMessages(ssems []interface{}) []*ServerSdkEventMessage {
	snams := make([]*ServerSdkEventMessage, 0)
	for index, _ := range ssems {
		ssem := ssems[index].(*ServerSdkEventMessage)
		snams = append(snams, CreateSaasNativeAppMessage(*ssem))
	}
	return snams
}
