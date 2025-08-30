package disk

import (
	"context"
	"errors"
	"goframe-ex/equeue/inter"
	"io"
	"os"
	"sync"
	"time"
)

const (
	filePerm  = 0600     // 数据写入权限
	indexFile = ".index" // 消息索引文件
)

type Queue struct {
	sync.RWMutex
	close  bool
	ticker *time.Ticker
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	writer *writer
	reader *reader
}

func New(config *inter.MqConfig) (queue *Queue, err error) {
	if _, err = os.Stat(config.Path); err != nil {
		return
	}
	queue = &Queue{close: false, wg: &sync.WaitGroup{}, writer: &writer{config: config}, reader: &reader{config: config}}
	queue.ticker = time.NewTicker(config.BatchTime)
	queue.ctx, queue.cancel = context.WithCancel(context.TODO())
	err = queue.reader.restore()
	if err != nil {
		return
	}
	go queue.sync()
	return
}

// Write data
func (q *Queue) Write(data []byte) error {
	if q.close {
		return errors.New("closed")
	}

	q.Lock()
	defer q.Unlock()

	return q.writer.write(data)
}

// Read data
func (q *Queue) Read() (int64, int64, []byte, error) {
	if q.close {
		return 0, 0, nil, errors.New("closed")
	}

	q.RLock()
	defer q.RUnlock()

	index, offset, data, err := q.reader.read()
	if err == io.EOF && (q.writer.file == nil || q.reader.file.Name() != q.writer.file.Name()) {
		_ = q.reader.safeRotate()
	}
	return index, offset, data, err
}

// Commit index and offset
func (q *Queue) Commit(index int64, offset int64) {
	if q.close {
		return
	}

	ck := &q.reader.checkpoint
	ck.Index, ck.Offset = index, offset
	q.reader.sync()
}

// Close Queue
func (q *Queue) Close() {
	if q.close {
		return
	}

	q.close = true
	q.cancel()
	q.wg.Wait()
	q.writer.close()
	q.reader.close()
}

// sync data
func (q *Queue) sync() {
	q.wg.Add(1)
	defer q.wg.Done()
	for {
		select {
		case <-q.ticker.C:
			q.Lock()
			q.writer.sync()
			q.Unlock()
		case <-q.ctx.Done():
			return
		}
	}
}
