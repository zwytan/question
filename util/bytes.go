package util

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

func ToBytes(s string) (uint64, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	//返回首字母所在位置
	i := strings.IndexFunc(s, unicode.IsLetter)

	if i == -1 {
		return 0, errors.New("参数中未包含单位容量")
	}

	n, b := s[:i], s[i:]
	bytes, err := strconv.ParseInt(n, 10, 64)
	if err != nil || bytes <= 0 {
		return 0, errors.New("参数错误")
	}

	switch b {
	case "GB":
		return uint64(bytes * GB), nil
	case "MB":
		return uint64(bytes * MB), nil
	case "KB":
		return uint64(bytes * KB), nil
	default:
		return 0, errors.New("参数错误")
	}
}
