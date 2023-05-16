package redis_utils

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// type Conn interface {
// 	// Close closes the connection.
// 	Close() error
// 	// Err returns a non-nil value when the connection is not usable.
// 	Err() error
// 	// Do sends a command to the server and returns the received reply.
// 	// This function will use the timeout which was set when the connection is created
// 	Do(commandName string, args ...interface{}) (reply interface{}, err error)
// 	// Send writes the command to the client's output buffer.
// 	Send(commandName string, args ...interface{}) error
// 	// Flush flushes the output buffer to the Redis server.
// 	Flush() error
// 	// Receive receives a single reply from the Redis server
// 	Receive() (reply interface{}, err error)
// }

// opentelemetry
type OtelRedis struct {
	conn redis.Conn
}

func (m OtelRedis) Close() error {
	return m.conn.Close()
}

func (m OtelRedis) Err() error {
	return m.conn.Err()
}

func (m OtelRedis) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	
	return m.conn.Do(commandName, args)
}

func (m OtelRedis) Send(commandName string, args ...interface{}) error {
	fmt.Println("send start")
	err := m.conn.Send(commandName, args)
	fmt.Println("send end")
	return err
}

func (m OtelRedis) Flush() error {
	return m.conn.Flush()
}

func (m OtelRedis) Receive() (reply interface{}, err error) {
	return m.conn.Receive()
}
