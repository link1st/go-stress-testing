// Package model 数据模型
package model

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mattn/go-shellwords"

	"github.com/link1st/go-stress-testing/helper"
)

// CURL curl参数解析
type CURL struct {
	Data map[string][]string
}

// getDataValue 获取数据
func (c *CURL) getDataValue(keys []string) []string {
	var (
		value = make([]string, 0)
	)
	for _, key := range keys {
		var (
			ok bool
		)
		value, ok = c.Data[key]
		if ok {
			break
		}
	}
	return value
}

// ParseTheFile 从文件中解析curl
func ParseTheFile(path string) (curl *CURL, err error) {
	if path == "" {
		err = errors.New("路径不能为空")
		return
	}
	curl = &CURL{
		Data: make(map[string][]string),
	}
	file, err := os.Open(path)
	if err != nil {
		err = errors.New("打开文件失败:" + err.Error())
		return
	}
	defer func() {
		_ = file.Close()
	}()
	dataBytes, err := io.ReadAll(file)
	if err != nil {
		err = errors.New("读取文件失败:" + err.Error())
		return
	}
	args, err := shellwords.Parse(string(dataBytes))
	if err != nil {
		err = errors.New("解析文件失败:" + err.Error())
		return
	}
	args = argsTrim(args)
	var key string
	for _, arg := range args {
		arg = removeSpaces(arg)
		if arg == "" {
			continue
		}
		if isURL(arg) {
			curl.Data[keyCurl] = append(curl.Data[keyCurl], arg)
			key = ""
			continue
		}
		if isKey(arg) {
			key = arg
			continue
		}
		curl.Data[key] = append(curl.Data[key], arg)
	}
	return
}

func argsTrim(args []string) []string {
	result := make([]string, 0)
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if arg == "\n" {
			continue
		}
		if strings.Contains(arg, "\n") {
			arg = strings.ReplaceAll(arg, "\n", "")
		}
		if strings.Index(arg, "-X") == 0 {
			result = append(result, arg[0:2])
			result = append(result, arg[2:])
		} else {
			result = append(result, arg)
		}
	}
	return result
}

func removeSpaces(data string) string {
	data = strings.TrimFunc(data, func(r rune) bool {
		if r == ' ' || r == '\\' || r == '\n' {
			return true
		}
		return false
	})
	return data
}

func isKey(data string) bool {
	return strings.HasPrefix(data, "-") || strings.HasPrefix(data, keyCurl)
}

func isURL(data string) bool {
	return strings.HasPrefix(data, "http://") || strings.HasPrefix(data, "https://")
}

// String string
func (c *CURL) String() (url string) {
	curlByte, _ := json.Marshal(c)
	return string(curlByte)
}

const (
	keyCurl = "curl"
)

// GetURL 获取url
func (c *CURL) GetURL() (url string) {
	keys := []string{keyCurl, "--url", "--location"}
	value := c.getDataValue(keys)
	if len(value) <= 0 {
		return
	}
	url = value[0]
	return
}

// GetMethod 获取 请求方式
func (c *CURL) GetMethod() (method string) {
	keys := []string{"-X", "--request"}
	value := c.getDataValue(keys)
	if len(value) <= 0 {
		return c.defaultMethod()
	}
	method = strings.ToUpper(value[0])
	if helper.InArrayStr(method, []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}) {
		return method
	}
	return c.defaultMethod()
}

// defaultMethod 获取默认方法
func (c *CURL) defaultMethod() (method string) {
	method = http.MethodGet
	body := c.GetBody()
	if len(body) > 0 {
		return http.MethodPost
	}
	return
}

// GetHeaders 获取请求头
func (c *CURL) GetHeaders() (headers map[string]string) {
	headers = make(map[string]string)
	keys := []string{"-H", "--header"}
	value := c.getDataValue(keys)
	for _, v := range value {
		getHeaderValue(v, headers)
	}
	return
}

// GetHeadersStr 获取请求头string
func (c *CURL) GetHeadersStr() string {
	headers := c.GetHeaders()
	bytes, _ := json.Marshal(&headers)
	return string(bytes)
}

// GetBody 获取body
// 返回字符串格式body，不支持二进制文件
func (c *CURL) GetBody() (body string) {
	keys := []string{"--data", "-d", "--data-urlencode", "--data-raw", "--data-binary"}
	value := c.getDataValue(keys)
	if len(value) <= 0 {
		body = c.getPostForm()
		return
	}
	body = value[0]
	return
}

// GetBodyBytes 获取body二进制数据
// 支持 --data、-d、--data-binary 等参数的 @filename 语法，从文件读取内容
// 注意：--data-raw 不支持 @filename（curl 原生行为）
func (c *CURL) GetBodyBytes() (bodyBytes []byte, err error) {
	// 处理所有支持 @filename 的 data 参数
	// 注意：--data-raw 故意不包含在内，因为它不处理 @ 符号
	keys := []string{"--data", "-d", "--data-urlencode", "--data-binary"}
	value := c.getDataValue(keys)
	if len(value) <= 0 {
		return nil, nil
	}

	data := value[0]
	// 检查是否以 @ 开头，表示从文件读取
	if strings.HasPrefix(data, "@") {
		filePath := strings.TrimPrefix(data, "@")
		// 读取文件内容
		bodyBytes, err = os.ReadFile(filePath)
		if err != nil {
			return nil, errors.New("读取文件失败:" + err.Error())
		}
		return bodyBytes, nil
	}

	// 如果不是文件路径，直接返回字符串的字节数组
	return []byte(data), nil
}

// getPostForm get post form
func (c *CURL) getPostForm() (body string) {
	keys := []string{"--form", "-F", "--form-string"}
	value := c.getDataValue(keys)
	if len(value) <= 0 {
		return
	}
	body = strings.Join(value, "&")
	return
}
