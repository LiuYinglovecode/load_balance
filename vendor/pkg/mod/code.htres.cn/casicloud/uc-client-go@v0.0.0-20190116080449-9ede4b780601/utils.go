package uclient //import "code.htres.cn/casicloud/uc-client-go"

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type keyPair struct {
	key  string
	pair string
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func signQueryMap(params map[string]string, clientSecret string) (sign string, err error) {
	var buffer bytes.Buffer
	// sort by key
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for i := range keys {
		buffer.WriteString(keys[i])
		buffer.WriteString("=")
		buffer.WriteString(params[keys[i]])
	}

	if len(clientSecret) > 0 {
		buffer.WriteString(clientSecret)
	}

	// src := buffer.String()
	// fmt.Println(src)
	md5str := fmt.Sprintf("%x", md5.Sum(buffer.Bytes()))
	return md5str, nil
}

func queryToMap(query string) map[string]string {
	m := make(map[string]string)
	if len(query) == 0 {
		return m
	}

	s := strings.Split(query, "&")
	for i := range s {
		k := s[i]
		p := strings.Split(k, "=")
		if len(p) == 2 {
			m[p[0]] = p[1]
		}
	}

	return m
}

func signStringQuery(query string, clientSecret string) (signedQuery string, err error) {
	var buffer bytes.Buffer
	if len(query) > 0 {
		buffer.WriteString(query)
	}
	if !strings.Contains(query, "ts=") {
		buffer.WriteString("&ts=")
		buffer.WriteString(strconv.FormatInt(makeTimestamp(), 10))
	}

	m := queryToMap(buffer.String())
	sign, err := signQueryMap(m, clientSecret)
	if err != nil {
		return query, err
	}
	buffer.WriteString("&sign=")
	buffer.WriteString(sign)
	return buffer.String(), nil
}
