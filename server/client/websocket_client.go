/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-15
* Time: 21:03
 */

package client

import (
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"net/url"
	"strings"
)

const (
	connRetry = 3 // 建立连接重试次数
)

type WebSocket struct {
	conn    *websocket.Conn
	UrlLink string
	Url     *url.URL
	IsSsl   bool
}

func NewWebSocket(urlLink string) (ws *WebSocket) {
	var (
		isSsl bool
	)

	if strings.HasPrefix(urlLink, "wss://") {
		isSsl = true
	}

	u, err := url.Parse(urlLink)
	// 解析失败
	if err != nil {
		panic(err)
	}

	ws = &WebSocket{
		UrlLink: urlLink,
		Url:     u,
		IsSsl:   isSsl,
	}
	return
}

func (w *WebSocket) getLink() (link string) {
	link = w.UrlLink

	return
}

func (w *WebSocket) getOrigin() (origin string) {
	origin = "http://"
	if w.IsSsl {
		origin = "https://"
	}

	origin = fmt.Sprintf("%s%s/", origin, w.Url.Host)

	return
}

// 关闭
func (w *WebSocket) Close() (err error) {
	if w == nil {

		return
	}

	if w.conn == nil {
		return
	}

	w.conn.Close()

	return
}

func (w *WebSocket) GetConn() (err error) {

	var (
		conn *websocket.Conn
		i    int
	)

	for i = 0; i < connRetry; i++ {
		conn, err = websocket.Dial(w.getLink(), "", w.getOrigin())
		if err != nil {
			fmt.Println("GetConn 建立连接失败 in...", i, err)

			continue
		}
		w.conn = conn

		return
	}

	if err != nil {
		fmt.Println("GetConn 建立连接失败", i, err)
	}

	return
}

// 发送数据
func (w *WebSocket) Write(body []byte) (err error) {
	if w.conn == nil {
		err = errors.New("未建立连接")

		return
	}

	_, err = w.conn.Write(body)
	if err != nil {
		fmt.Println("发送数据失败:", err)

		return
	}

	return
}

// 接收数据
func (w *WebSocket) Read() (msg []byte, err error) {
	if w.conn == nil {
		err = errors.New("未建立连接")

		return
	}

	msg = make([]byte, 512)

	n, err := w.conn.Read(msg)
	if err != nil {
		fmt.Println("接收数据失败:", err)

		return nil, err
	}

	return msg[:n], nil
}
