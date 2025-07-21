package adapter

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
	"time"
)

// 服务处理获取前缀
type AdapterRedis struct {
	redis *gredis.Redis
	pre   string
}

// NewAdapterRedis creates and returns a new memory cache object.
func NewAdapterRedis(redis *gredis.Redis, name ...string) *AdapterRedis {
	var p = ""
	if len(name) > 0 {
		p = name[0]
	}
	if p != "" && !strings.HasSuffix(p, ":") {
		p = p + ":"
	}

	return &AdapterRedis{
		redis: redis,
		pre:   p,
	}
}

func (c *AdapterRedis) checkKey(key interface{}) string {
	redisKey := gconv.String(key)
	if !strings.HasPrefix(redisKey, c.pre) {
		redisKey = c.pre
	}
	return redisKey
}

func (c *AdapterRedis) Set(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (err error) {
	redisKey := c.checkKey(key)
	if value == nil || duration < 0 {
		_, err = c.redis.Del(ctx, redisKey)
	} else {
		if duration == 0 {
			_, err = c.redis.Set(ctx, redisKey, value)
		} else {
			_, err = c.redis.Set(ctx, redisKey, value, gredis.SetOption{TTLOption: gredis.TTLOption{PX: gconv.PtrInt64(duration.Milliseconds())}})
		}
	}
	return err
}

func (c *AdapterRedis) SetMap(ctx context.Context, mapdata map[interface{}]interface{}, duration time.Duration) error {
	if len(mapdata) == 0 {
		return nil
	}
	data := make(map[string]interface{}, len(mapdata))
	for k, v := range data {
		data[c.checkKey(k)] = v
	}

	// DEL.
	if duration < 0 {
		var (
			index = 0
			keys  = make([]string, len(data))
		)
		for k := range data {
			keys[index] = k
			index += 1
		}
		_, err := c.redis.Del(ctx, keys...)
		if err != nil {
			return err
		}
	}
	if duration == 0 {
		err := c.redis.MSet(ctx, data)
		if err != nil {
			return err
		}
	}
	if duration > 0 {
		var err error
		for k, v := range data {
			if err = c.Set(ctx, k, v, duration); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *AdapterRedis) SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (bool, error) {
	var (
		err      error
		redisKey = c.checkKey(key)
	)
	// Execute the function and retrieve the result.
	f, ok := value.(Func)
	if !ok {
		// Compatible with raw function value.
		f, ok = value.(func(ctx context.Context) (value interface{}, err error))
	}
	if ok {
		if value, err = f(ctx); err != nil {
			return false, err
		}
	}
	// DEL.
	if duration < 0 || value == nil {
		var delResult int64
		delResult, err = c.redis.Del(ctx, redisKey)
		if err != nil {
			return false, err
		}
		if delResult == 1 {
			return true, err
		}
		return false, err
	}
	ok, err = c.redis.SetNX(ctx, redisKey, value)
	if err != nil {
		return ok, err
	}
	if ok && duration > 0 {
		// Set the expiration.
		_, err = c.redis.PExpire(ctx, redisKey, duration.Milliseconds())
		if err != nil {
			return ok, err
		}
		return ok, err
	}
	return ok, err
}

func (c *AdapterRedis) SetIfNotExistFunc(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration) (ok bool, err error) {
	value, err := f(ctx)
	if err != nil {
		return false, err
	}
	return c.SetIfNotExist(ctx, key, value, duration)
}

func (c *AdapterRedis) SetIfNotExistFuncLock(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration) (ok bool, err error) {
	value, err := f(ctx)
	if err != nil {
		return false, err
	}
	return c.SetIfNotExist(ctx, key, value, duration)
}

func (c *AdapterRedis) Get(ctx context.Context, key interface{}) (*gvar.Var, error) {
	return c.redis.Get(ctx, c.checkKey(key))
}

func (c *AdapterRedis) GetOrSet(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (result *gvar.Var, err error) {
	result, err = c.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if result.IsNil() {
		return gvar.New(value), c.Set(ctx, key, value, duration)
	}
	return
}

func (c *AdapterRedis) GetOrSetFunc(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration) (result *gvar.Var, err error) {
	v, err := c.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if v.IsNil() {
		value, err := f(ctx)
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, nil
		}
		return gvar.New(value), c.Set(ctx, key, value, duration)
	} else {
		return v, nil
	}
}

func (c *AdapterRedis) GetOrSetFuncLock(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration) (result *gvar.Var, err error) {
	return c.GetOrSetFunc(ctx, key, f, duration)
}

func (c *AdapterRedis) Contains(ctx context.Context, key interface{}) (bool, error) {
	n, err := c.redis.Exists(ctx, c.checkKey(key))
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (c *AdapterRedis) Size(ctx context.Context) (size int, err error) {
	keys, err := c.redis.Keys(ctx, c.pre+"*")
	//if err != nil {
	//	return 0, err
	//}
	return len(keys), err
}

func (c *AdapterRedis) Data(ctx context.Context) (map[interface{}]interface{}, error) {
	// Keys.
	keys, err := c.redis.Keys(ctx, c.pre+"*")
	if err != nil {
		return nil, err
	}
	// Key-Value pairs.
	var m map[string]*gvar.Var
	m, err = c.redis.MGet(ctx, keys...)
	if err != nil {
		return nil, err
	}
	// Type converting.
	data := make(map[interface{}]interface{})
	for k, v := range m {
		data[k] = v.Val()
	}
	return data, nil
}

func (c *AdapterRedis) Keys(ctx context.Context) ([]interface{}, error) {
	keys, err := c.redis.Keys(ctx, c.pre+"*")
	if err != nil {
		return nil, err
	}
	return gconv.Interfaces(keys), nil
}

func (c *AdapterRedis) Values(ctx context.Context) (values []interface{}, err error) {
	// Keys.
	keys, err := c.redis.Keys(ctx, c.pre+"*")
	if err != nil {
		return nil, err
	}
	// Key-Value pairs.
	var m map[string]*gvar.Var
	m, err = c.redis.MGet(ctx, keys...)
	if err != nil {
		return nil, err
	}
	// Values.
	for _, key := range keys {
		if v := m[key]; !v.IsNil() {
			values = append(values, v.Val())
		}
	}
	return values, nil
}

func (c *AdapterRedis) Update(ctx context.Context, key interface{}, value interface{}) (oldValue *gvar.Var, exist bool, err error) {
	var (
		v        *gvar.Var
		oldPTTL  int64
		redisKey = c.checkKey(key)
	)
	// TTL.
	oldPTTL, err = c.redis.PTTL(ctx, redisKey) // update ttl -> pttl(millisecond)
	if err != nil {
		return
	}
	if oldPTTL == -2 || oldPTTL == 0 {
		// It does not exist or expired.
		return
	}
	// Check existence.
	v, err = c.redis.Get(ctx, redisKey)
	if err != nil {
		return
	}
	oldValue = v
	// DEL.
	if value == nil {
		_, err = c.redis.Del(ctx, redisKey)
		if err != nil {
			return
		}
		return
	}
	// Update the value.
	if oldPTTL == -1 {
		_, err = c.redis.Set(ctx, redisKey, value)
	} else {
		// update SetEX -> SET PX Option(millisecond)
		// Starting with Redis version 2.6.12: Added the EX, PX, NX and XX options.
		_, err = c.redis.Set(ctx, redisKey, value, gredis.SetOption{TTLOption: gredis.TTLOption{PX: gconv.PtrInt64(oldPTTL)}})
	}
	return oldValue, true, err
}

func (c *AdapterRedis) UpdateExpire(ctx context.Context, key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	var (
		v        *gvar.Var
		oldPTTL  int64
		redisKey = c.checkKey(key)
	)
	// TTL.
	oldPTTL, err = c.redis.PTTL(ctx, redisKey)
	if err != nil {
		return
	}
	if oldPTTL == -2 || oldPTTL == 0 {
		// It does not exist or expired.
		oldPTTL = -1
		return
	}
	oldDuration = time.Duration(oldPTTL) * time.Millisecond
	// DEL.
	if duration < 0 {
		_, err = c.redis.Del(ctx, redisKey)
		return
	}
	// Update the expiration.
	if duration > 0 {
		_, err = c.redis.PExpire(ctx, redisKey, duration.Milliseconds())
	}
	// No expire.
	if duration == 0 {
		v, err = c.redis.Get(ctx, redisKey)
		if err != nil {
			return
		}
		_, err = c.redis.Set(ctx, redisKey, v.Val())
	}
	return
}

func (c AdapterRedis) GetExpire(ctx context.Context, key interface{}) (time.Duration, error) {
	pttl, err := c.redis.PTTL(ctx, c.checkKey(key))
	if err != nil {
		return 0, err
	}
	switch pttl {
	case -1:
		return 0, nil
	case -2, 0: // It does not exist or expired.
		return -1, nil
	default:
		return time.Duration(pttl) * time.Millisecond, nil
	}
}

func (c *AdapterRedis) Remove(ctx context.Context, keys ...interface{}) (lastValue *gvar.Var, err error) {
	if len(keys) == 0 {
		return nil, nil
	}
	// Retrieves the last key value.
	if lastValue, err = c.redis.Get(ctx, c.checkKey(keys[len(keys)-1])); err != nil {
		return nil, err
	}
	var keystr = make([]string, len(keys))
	for i, key := range keys {
		keystr[i] = c.checkKey(key)
	}

	// Deletes all given keys.
	_, err = c.redis.Del(ctx, keystr...)
	return
}

func (c *AdapterRedis) Clear(ctx context.Context) error {
	keys, err := c.Keys(ctx)
	if err != nil {
		return err
	}
	_, err = c.Remove(ctx, keys...)

	return err
}

func (c *AdapterRedis) Close(ctx context.Context) error {
	return nil
}
