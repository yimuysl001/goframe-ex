package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"goframe-ex/equeue/inter"
	"math/rand"
	"time"
)

// SendByteMsg 生产数据
func (r *RedisMq) SendByteMsg(ctx context.Context, topic string, body []byte) (mqMsg inter.MqMsg, err error) {
	//if r.poolName == "" {
	//	return mqMsg, gerror.New("RedisMq producer not register")
	//}

	if topic == "" {
		return mqMsg, gerror.New("RedisMq topic is empty")
	}

	mqMsg = inter.MqMsg{
		RunType:   inter.SendMsg,
		Topic:     topic,
		MsgId:     getRandMsgId(),
		Body:      body,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(mqMsg)
	if err != nil {
		return
	}

	key := r.genKey(r.groupName, topic)
	if _, err = g.Redis(r.poolName).Do(ctx, "LPUSH", key, data); err != nil {
		return
	}

	if r.timeout > 0 {
		if _, err = g.Redis(r.poolName).Do(ctx, "EXPIRE", key, time.Duration(r.timeout)*time.Second); err != nil {
			return
		}
	}

	return
}

// SendMsg 按字符串类型生产数据
func (r *RedisMq) SendMsg(ctx context.Context, topic string, body string) (mqMsg inter.MqMsg, err error) {
	return r.SendByteMsg(ctx, topic, []byte(body))
}

func (r *RedisMq) SendDelayMsg(ctx context.Context, topic string, body string, delaySecond int64) (mqMsg inter.MqMsg, err error) {
	if delaySecond < 1 {
		return r.SendMsg(ctx, topic, body)
	}

	if r.poolName == "" {
		err = gerror.New("SendDelayMsg RedisMq not register")
		return
	}

	if topic == "" {
		err = gerror.New("SendDelayMsg RedisMq topic is empty")
		return
	}

	mqMsg = inter.MqMsg{
		RunType:   inter.SendMsg,
		Topic:     topic,
		MsgId:     getRandMsgId(),
		Body:      []byte(body),
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(mqMsg)
	if err != nil {
		return
	}

	var (
		conn         = g.Redis(r.poolName)
		key          = r.genKey(r.groupName, "delay:"+topic)
		expireSecond = time.Now().Unix() + delaySecond
		timePiece    = fmt.Sprintf("%s:%d", key, expireSecond)
		z            = gredis.ZAddMember{Score: float64(expireSecond), Member: timePiece}
	)

	if _, err = conn.ZAdd(ctx, key, &gredis.ZAddOption{}, z); err != nil {
		return
	}

	if _, err = conn.RPush(ctx, timePiece, data); err != nil {
		return
	}

	// consumer will also delete the item
	if r.timeout > 0 {
		_, _ = conn.Expire(ctx, timePiece, int64(time.Duration(r.timeout)*time.Second)+delaySecond)
		_, _ = conn.Expire(ctx, key, int64(time.Duration(r.timeout)*time.Second))
	}

	return
}

func getRandMsgId() string {
	rand.NewSource(time.Now().UnixNano())
	radium := rand.Intn(999) + 1
	timeCode := time.Now().UnixNano()
	return fmt.Sprintf("%d%.4d", timeCode, radium)
}

// 生成队列key
func (r *RedisMq) genKey(groupName string, topic string) string {
	return fmt.Sprintf("queue:%s_%s", groupName, topic)
}
