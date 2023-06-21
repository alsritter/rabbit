package event

import (
	"context"

	bus "github.com/alsritter/nsq-event-bus"
)

type Message struct {
	MessageID string
	Timestamp int64
	Body      any
}

type Option func(*Server)

func WithHandlerConcurrency(concurrency int) Option {
	return func(s *Server) {
		s.handlerConcurrency = concurrency
	}
}

func WithMaxAttempts(attempts int) Option {
	return func(s *Server) {
		s.maxAttempts = attempts
	}
}

type HandlerFunc func(message *Message) error

type Subscriber interface {
	HandleMessage(message *Message) error // 处理消息
	GetName() string                      // 取得观察者的名称，即同个观察者只能消费一次
	GetEventName() string                 // 取得观察者订阅的事件名称
}

type Server struct {
	handlerConcurrency int // 消费者并发数
	maxAttempts        int // 最大尝试次数
	subscriberList     []Subscriber
}

func NewEvnetServer(options ...Option) *Server {
	srv := &Server{
		handlerConcurrency: 4,
		maxAttempts:        5,
	}

	for _, option := range options {
		option(srv)
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	for _, subscriber := range s.subscriberList {
		if err := bus.On(bus.ListenerConfig{
			Topic:              subscriber.GetEventName(),
			Channel:            subscriber.GetName(),
			HandlerFunc:        s.WrapHandler(subscriber.HandleMessage),
			HandlerConcurrency: s.handlerConcurrency,
			MaxAttempts:        uint16(s.maxAttempts),
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return nil
}

func (s *Server) Register(subscribers ...Subscriber) {
	s.subscriberList = append(s.subscriberList, subscribers...)
}

func (s *Server) WrapHandler(handler HandlerFunc) bus.HandlerFunc {

	return nil
}
