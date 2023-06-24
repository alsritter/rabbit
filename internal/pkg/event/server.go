package event

import (
	"context"

	bus "github.com/alsritter/nsq-event-bus"
)

var ServerInstance *Server

type Message struct {
	MessageID string
	Timestamp int64
	Body      string
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

func WithMaxInFlight(inFlight int) Option {
	return func(s *Server) {
		s.maxInFlight = inFlight
	}
}

type HandlerFunc func(message *Message) error

type Subscriber interface {
	HandleMessage(message *Message) error // 处理消息
	GetName() string                      // 取得观察者的名称，即同个观察者只能消费一次
	GetEventName() string                 // 取得观察者订阅的事件名称
}

type Server struct {
	lookup             []string // nsqlookupd 地址
	address            string   // nsqd 地址
	maxInFlight        int      // 最大并发数
	handlerConcurrency int      // 消费者并发数
	maxAttempts        int      // 最大尝试次数
	subscriberList     []Subscriber

	emitter *bus.Emitter // 发布者
}

func NewEventServer(address string, lookup []string, options ...Option) *Server {
	srv := &Server{
		lookup:             lookup,
		address:            address,
		handlerConcurrency: 4,
		maxAttempts:        5,
	}

	for _, option := range options {
		option(srv)
	}

	var err error
	srv.emitter, err = bus.NewEmitter(bus.EmitterConfig{
		Address:     srv.address,
		MaxInFlight: srv.maxInFlight,
	})
	if err != nil {
		panic(err)
	}

	ServerInstance = srv
	return srv
}

func (s *Server) Start(ctx context.Context) error {
	for _, subscriber := range s.subscriberList {
		if err := bus.On(bus.ListenerConfig{
			Lookup:             s.lookup,
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
	return bus.HandlerFunc(func(message *bus.Message) (any, error) {
		if message.Attempts > uint16(s.maxAttempts) {
			message.Finish()
			return nil, nil
		}

		m := Message{
			MessageID: string(message.ID[:]),
			Timestamp: message.Timestamp,
			Body:      string(message.Payload),
		}

		if err := handler(&m); err != nil {
			message.Requeue(-1)
			return nil, err
		}

		message.Finish()
		return nil, nil
	})
}

/**
 * 事件发布，入参是指针
 */
func (s *Server) Emit(eventName string, payload any) error {
	return s.emitter.Emit(eventName, payload)
}

/**
 * 事件发布，入参是指针
 */
func (s *Server) EmitAsync(eventName string, payload any) error {
	return s.emitter.EmitAsync(eventName, payload)
}
