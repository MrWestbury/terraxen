package services

import (
	"fmt"
	"net/url"
	"strings"
)

type Options struct {
	Hostname  string
	Username  string
	Password  string
	Database  string
	ExtraOpts map[string]string
}

func (opts Options) ConnectionString() string {
	connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s", url.QueryEscape(opts.Username), url.QueryEscape(opts.Password), opts.Hostname, opts.Database)
	var extraElems []string
	for k, v := range opts.ExtraOpts {
		key := url.QueryEscape(k)
		value := url.QueryEscape(v)
		extraElemStr := fmt.Sprintf("%s=%s", key, value)
		extraElems = append(extraElems, extraElemStr)
	}
	if len(extraElems) > 0 {
		extraStr := strings.Join(extraElems, "&")
		connStr = fmt.Sprintf("%s?%s", connStr, extraStr)
	}
	return connStr
}
