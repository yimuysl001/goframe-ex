package disk

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"goframe-ex/equeue/driver/logger"
	"goframe-ex/equeue/inter"
	"sync"
	"time"
)

type DiskProducerMq struct {
	config    inter.MqConfig
	producers map[string]*Queue
	sync.Mutex
}

type DiskConsumerMq struct {
	config inter.MqConfig
	queue  map[string]*Queue
	flag   map[string]bool
}

func NewDiskQueue(topic string, config inter.MqConfig) *Queue {
	var conf = new(inter.MqConfig)
	err := gjson.New(config).Scan(&conf)
	if err != nil {
		panic(err)
	}

	conf.Path = fmt.Sprintf(config.Path + "/" + config.GroupName + "/" + topic)
	conf.BatchTime = config.BatchTime * time.Second

	if !gfile.Exists(conf.Path) {
		if err := gfile.Mkdir(conf.Path); err != nil {
			logger.Logger().Errorf(gctx.GetInitCtx(), "NewDiskQueue Failed to create the cache directory. Procedure, err:%+v", err)
			return nil
		}
	}

	queue, err := New(conf)
	if err != nil {
		logger.Logger().Errorf(gctx.GetInitCtx(), "NewDiskQueue err:%v", err)
		return nil
	}
	return queue
}
