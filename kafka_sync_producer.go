/**
 *	Copyright 2020 Beijing Volcano Engine Technology Co., Ltd.
 *	Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *  http://www.apache.org/licenses/LICENSE-2.0
 *	Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

package datarangers_sdk

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

type SyncKafkaProducer struct {
	syncProducer sarama.SyncProducer
	config       *KafkaProducerConf
	topic        string
}

type KeyPair struct {
	Key     string
	Body    []byte
	Headers []sarama.RecordHeader
}

func GetSaramaProducerConfig(globalKafkaConf *KafkaConfig) *sarama.Config {
	sarama.Logger = infoLog
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	if globalKafkaConf.KafkaProducerConf != nil {
		if globalKafkaConf.KafkaProducerConf.RequiredAcks == 0 {
			config.Producer.RequiredAcks = sarama.WaitForLocal
		} else {
			config.Producer.RequiredAcks = sarama.RequiredAcks(globalKafkaConf.KafkaProducerConf.RequiredAcks)
		}

		if globalKafkaConf.KafkaProducerConf.RetryConfig == nil {
			config.Producer.Retry.Max = 3
		} else {
			config.Producer.Retry.Max = globalKafkaConf.KafkaProducerConf.RetryConfig.Max
		}

		//flush配置
		if globalKafkaConf.KafkaProducerConf.FlushConfig != nil {
			config.Producer.Flush.Frequency = time.Duration(globalKafkaConf.KafkaProducerConf.FlushConfig.FrequencyMs) * time.Millisecond
			config.Producer.Flush.Messages = globalKafkaConf.KafkaProducerConf.FlushConfig.Messages
			config.Producer.Flush.MaxMessages = globalKafkaConf.KafkaProducerConf.FlushConfig.MaxMessages
			config.Producer.Flush.Bytes = globalKafkaConf.KafkaProducerConf.FlushConfig.Bytes
		}
	}

	config.Producer.Return.Successes = true

	if globalKafkaConf.KafkaSaslType == 1 {
		// plain text
		config.Net.SASL.Enable = true
		config.Net.SASL.User = globalKafkaConf.KafkaSaslUser
	} else if globalKafkaConf.KafkaSaslType == 2 {
		// kerberos
		config.Net.SASL.Enable = true
		config.Net.SASL.GSSAPI.DisablePAFXFAST = true
		config.Net.SASL.Mechanism = sarama.SASLTypeGSSAPI
		config.Net.SASL.GSSAPI.AuthType = sarama.KRB5_KEYTAB_AUTH
		config.Net.SASL.GSSAPI.KerberosConfigPath = globalKafkaConf.GssApiKerberosConfigPath
		config.Net.SASL.GSSAPI.KeyTabPath = globalKafkaConf.GssApiKeyTabPath

		if globalKafkaConf.GssApiPrincipal != "" {
			splits := strings.Split(globalKafkaConf.GssApiPrincipal, "@")
			if len(splits) == 2 {
				config.Net.SASL.GSSAPI.Username = splits[0]
				config.Net.SASL.GSSAPI.Realm = splits[1]
				config.Net.SASL.GSSAPI.ServiceName = "kafka"
			}
		} else {
			config.Net.SASL.GSSAPI.Username = globalKafkaConf.KafkaSaslUser
			config.Net.SASL.GSSAPI.Realm = globalKafkaConf.GssApiRealm
			config.Net.SASL.GSSAPI.ServiceName = globalKafkaConf.GssApiServiceName
		}
	}

	return config
}

func NewSyncKafkaProducer(globalKafkaConf *KafkaConfig) (*SyncKafkaProducer, error) {
	if globalKafkaConf.Topic == "" {
		globalKafkaConf.Topic = KAFAKA_TOPIC
	}
	syncProducer, err := sarama.NewSyncProducer(globalKafkaConf.KafkaBrokers, GetSaramaProducerConfig(globalKafkaConf))
	if err != nil {
		fatal("NewSyncProducer Error: " + err.Error())
		return nil, err
	}

	return &SyncKafkaProducer{
		syncProducer: syncProducer,
		config:       globalKafkaConf.KafkaProducerConf,
		topic:        globalKafkaConf.Topic,
	}, nil
}

func (p *SyncKafkaProducer) Send(msgValue interface{}) error {
	marshal, err := json.Marshal(msgValue)
	if err != nil {
		fatal(err.Error())
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Value:     sarama.StringEncoder(marshal),
		Timestamp: time.Now(),
	}
	_, _, err = p.syncProducer.SendMessage(msg)
	if err != nil {
		fatal(err.Error())
		return err
	}
	return nil
}

func (p *SyncKafkaProducer) BatchSend(msgValues []interface{}) error {
	msgs := make([]*sarama.ProducerMessage, 0)
	for _, msgValue := range msgValues {
		marshal, err := json.Marshal(msgValue)
		if err != nil {
			fatal(err.Error())
			return err
		}
		msg := &sarama.ProducerMessage{
			Topic:     p.topic,
			Value:     sarama.StringEncoder(marshal),
			Timestamp: time.Now(),
		}
		// 暂时不用指定key
		//if keyPair.Key != "" {
		//	msg.Key = sarama.StringEncoder(keyPair.Key)
		//}

		msgs = append(msgs, msg)
	}
	return p.syncProducer.SendMessages(msgs)
}

func (p *SyncKafkaProducer) Close() error {
	return p.syncProducer.Close()
}
