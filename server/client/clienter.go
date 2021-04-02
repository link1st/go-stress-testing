// Package client clientpackage client
package client

// Clienter 接口 注册、连接、发送 等
type Clienter interface {
	GetConn() (err error)
	Close() (err error)
	Send()
}
