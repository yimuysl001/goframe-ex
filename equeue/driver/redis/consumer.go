package redis

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"goframe-ex/equeue/driver/logger"
	"goframe-ex/equeue/inter"
	"strconv"
	"time"
)

func (r *RedisMq) Unsubscribe(ctx context.Context, topic string) error {
	r.flag[topic] = true
	return nil
}

// ListenReceiveMsgDo 消费数据
func (r *RedisMq) ListenReceiveMsgDo(ctx context.Context, topic string, receiveDo func(ctx context.Context, mqMsg inter.MqMsg) error) (err error) {
	if r.poolName == "" {
		return gerror.New("RedisMq producer not register")
	}
	if topic == "" {
		return gerror.New("RedisMq topic is empty")
	}

	_, ok := r.flag[topic]
	if ok {
		return errors.New(topic + " is started ")
	}

	r.flag[topic] = false
	var (
		key  = r.genKey(r.groupName, topic)
		key2 = r.genKey(r.groupName, "delay:"+topic)
	)

	go func() {
		for range time.Tick(300 * time.Millisecond) {
			mqMsgList := r.loopReadQueue(ctx, topic, key)
			for _, mqMsg := range mqMsgList {
				err = receiveDo(ctx, mqMsg)
				if err != nil {
					logger.Logger().Warningf(ctx, "ListenReceiveMsgDo loopReadQueue err:%+v", err)
				}
			}
		}
	}()

	go func() {
		mqMsgCh, errCh := r.loopReadDelayQueue(ctx, topic, key2)
		for mqMsg := range mqMsgCh {
			err = receiveDo(ctx, mqMsg)
			if err != nil {
				logger.Logger().Warningf(ctx, "ListenReceiveMsgDo loopReadDelayQueue err:%+v", err)
			}
		}
		for err = range errCh {
			if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
				logger.Logger().Infof(ctx, "ListenReceiveMsgDo Delay topic:%v, err:%+v", topic, err)
			}
		}
	}()

	for {
		if r.flag[topic] {
			break
		}
		time.Sleep(1 * time.Second)

	}
	return nil
}

func (r *RedisMq) loopReadDelayQueue(ctx context.Context, topic, key string) (resCh chan inter.MqMsg, errCh chan error) {
	resCh = make(chan inter.MqMsg)
	errCh = make(chan error, 1)

	go func() {
		defer close(resCh)
		defer close(errCh)

		conn := g.Redis(r.poolName)
		for {
			if r.flag[topic] {
				break
			}
			now := time.Now().Unix()
			do, err := conn.Do(ctx, "zrangebyscore", key, "0", strconv.FormatInt(now, 10), "limit", 0, 1)
			if err != nil {
				return
			}
			val := do.Strings()
			if len(val) == 0 {
				select {
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				case <-time.After(time.Second):
					continue
				}
			}
			for _, listK := range val {
				for {
					pop, err := conn.LPop(ctx, listK)
					if err != nil {
						errCh <- err
						return
					} else if pop.IsEmpty() {
						_, _ = conn.ZRem(ctx, key, listK)
						_, _ = conn.Del(ctx, listK)
						break
					} else {
						var mqMsg inter.MqMsg
						if err = pop.Scan(&mqMsg); err != nil {
							g.Log().Warningf(ctx, "loopReadDelayQueue Scan err:%+v", err)
							break
						}

						if mqMsg.MsgId == "" {
							continue
						}

						select {
						case resCh <- mqMsg:
						case <-ctx.Done():
							errCh <- ctx.Err()
							return
						}
					}
				}
			}
		}
	}()
	return resCh, errCh
}

func (r *RedisMq) loopReadQueue(ctx context.Context, topic, key string) (mqMsgList []inter.MqMsg) {
	conn := g.Redis(r.poolName)
	for {
		if r.flag[topic] {
			break
		}
		data, err := conn.Do(ctx, "RPOP", key)
		if err != nil {
			logger.Logger().Warningf(ctx, "loopReadQueue redis RPOP err:%+v", err)
			break
		}

		if data.IsEmpty() {
			break
		}

		var mqMsg inter.MqMsg
		if err = data.Scan(&mqMsg); err != nil {
			logger.Logger().Warningf(ctx, "loopReadQueue Scan err:%+v", err)
			break
		}

		if mqMsg.MsgId != "" {
			mqMsgList = append(mqMsgList, mqMsg)
		}
	}
	return mqMsgList
}
