package inter

import "time"

const (
	_ = iota
	SendMsg
	ReceiveMsg
)

type MqConfig struct {
	Driver        string        `json:"driver"`
	Name          string        `json:"name"`
	Retry         int           `json:"retry"`
	Address       []string      `json:"address"`
	Version       string        `json:"version"`
	RandClient    bool          `json:"randClient"`
	MultiConsumer bool          `json:"multiConsumer"`
	Timeout       int64         `json:"timeout"`
	LogLevel      string        `json:"logLevel"`
	GroupName     string        `json:"groupName"`    //  组群名称
	Path          string        `json:"path"`         // disk 使用 数据存放路径
	BatchSize     int64         `json:"batchSize"`    // disk 使用  每N条消息同步一次，batchSize和batchTime满足其一就会同步一次
	BatchTime     time.Duration `json:"batchTime"`    // disk 使用  每N秒消息同步一次
	SegmentSize   int64         `json:"segmentSize"`  // disk 使用  每个topic分片数据文件最大字节
	SegmentLimit  int64         `json:"segmentLimit"` // disk 使用  每个topic最大分片数据文件数量
	UserName      string        `json:"userName"`
	Password      string        `json:"password"`
	Token         string        `json:"token"`
}
