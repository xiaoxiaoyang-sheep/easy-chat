/**
 * @author: Yanko/xiaoxiaoyang-sheep
 * @doc:
 **/

package websocket

import "net/http"

type DailOptions func(option *dailOption)

type dailOption struct {
	pattern string
	header  http.Header
}

func NewDailOptions(opts ...DailOptions) dailOption {
	o := dailOption{
		pattern: defaultPattern,
		header:  nil,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

func WithClientPattern(pattern string) DailOptions {
	return func(option *dailOption) {
		option.pattern = pattern
	}
}

func WithClientHeader(header http.Header) DailOptions {
	return func(option *dailOption) {
		option.header = header
	}
}
