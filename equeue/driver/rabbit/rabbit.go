package rabbit

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"goframe-ex/equeue/driver/logger"
	"goframe-ex/equeue/inter"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"math/rand"
	"sync"
)

type RabbitMQClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	config     *inter.MqConfig
	connected  bool
}

func NewRabbitMQClient(config inter.MqConfig) (*RabbitMQClient, error) {
	client := &RabbitMQClient{
		ctx:    gctx.New(),
		config: &config,
	}
	err := client.connect()
	if err != nil {
		return nil, err
	}
	if config.Kind == "" {
		config.Kind = "direct"
	}

	err = client.ensureExchange(config.Exchange, config.Kind)
	if err != nil {
		return nil, err
	}
	err = client.ensureQueue(config.Queue)

	if err != nil {
		return nil, err
	}
	err = client.bindQueue(config.Queue, config.Exchange, config.Key)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Connect 连接到RabbitMQ
func (c *RabbitMQClient) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var err error
	// 随机打乱节点顺序，实现简单负载均衡
	hosts := make([]string, len(c.config.Address))
	copy(hosts, c.config.Address)
	rand.Shuffle(len(hosts), func(i, j int) {
		hosts[i], hosts[j] = hosts[j], hosts[i]
	})

	if c.config.Name == "" {
		c.config.Name = "/"
	}
	// 尝试连接所有节点
	for _, host := range hosts {
		url := fmt.Sprintf("amqp://%s:%s@%s/%s",
			c.config.UserName, c.config.Password, host, c.config.Name)

		c.connection, err = amqp.Dial(url)
		if err == nil {
			break
		}
		logger.Logger().Error(c.ctx, "Failed to connect to %s: %v", host, err)
	}
	// 创建通道
	c.channel, err = c.connection.Channel()
	if err != nil {
		c.connection.Close()
		return fmt.Errorf("cannot create channel: %v", err)
	}

	// 监听连接关闭事件
	go c.monitorConnection()

	c.connected = true
	logger.Logger().Info(c.ctx, "Connected to RabbitMQ successfully")

	return err

}

// monitorConnection 监控连接状态
func (c *RabbitMQClient) monitorConnection() {
	closeChan := make(chan *amqp.Error)
	c.connection.NotifyClose(closeChan)

	select {
	case err := <-closeChan:
		if err != nil {
			logger.Logger().Error(c.ctx, "RabbitMQ connection closed: %v", err)
		}
		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()
		c.reconnect()
	case <-c.ctx.Done():
		return
	}
}

// reconnect 重连机制
func (c *RabbitMQClient) reconnect() {
	retryCount := 0
	for {
		if c.config.Retry > 0 && retryCount > c.config.Retry {
			logger.Logger().Error(c.ctx, "Reconnected failed")
			break
		}
		time.Sleep(time.Minute)
		if c.connected {
			continue
		}

		err := c.connect()
		if err == nil {
			logger.Logger().Info(c.ctx, "Reconnected to RabbitMQ successfully")
			return
		}

		logger.Logger().Error(c.ctx, "Reconnection failed: %v", err)
		retryCount++
	}
}

// EnsureExchange 声明交换机
func (c *RabbitMQClient) ensureExchange(name, kind string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return errors.New("not connected to RabbitMQ")
	}

	return c.channel.ExchangeDeclare(
		name,  // name
		kind,  // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

// EnsureQueue 声明队列
func (c *RabbitMQClient) ensureQueue(name string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return errors.New("not connected to RabbitMQ")
	}

	_, err := c.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}

// BindQueue 绑定队列到交换机
func (c *RabbitMQClient) bindQueue(queue, exchange, key string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.connected {
		return errors.New("not connected to RabbitMQ")
	}

	return c.channel.QueueBind(
		queue,    // queue name
		key,      // routing key
		exchange, // exchange
		false,    // no-wait
		nil,      // arguments
	)
}
