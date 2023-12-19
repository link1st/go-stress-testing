// Package client 客户端接口定义 client
package client

// Clienter 接口 获取连接、发送数据、关闭 等
type Clienter interface {
	GetConn() (err error)
	Send()
	Close() (err error)
}
