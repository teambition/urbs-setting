package tpl

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/service"
)

// Pagination 分页
type Pagination struct {
	PageToken string `json:"pageToken" query:"pageToken"`
	PageSize  int    `json:"pageSize,omitempty" query:"pageSize"`
	Skip      int    `json:"skip,omitempty" query:"skip"`
}

// Validate ...
func (pg *Pagination) Validate() error {
	if pg.Skip < 0 {
		pg.Skip = 0
	}

	if pg.PageSize > 1000 {
		return gear.ErrBadRequest.WithMsgf("pageSize(%v) should not great than 1000", pg.PageSize)
	}

	if pg.PageSize <= 0 {
		pg.PageSize = 10
	}

	return nil
}

// TokenToID 把 pageToken 转换为 int64
func (pg *Pagination) TokenToID() int64 {
	return PageTokenToID(pg.PageToken)
}

// TokenToTime 把 pageToken 转换为 time
func (pg *Pagination) TokenToTime() time.Time {
	return PageTokenToTime(pg.PageToken)
}

// PageTokenToID 把 pageToken 转换为 int64
func PageTokenToID(pageToken string) int64 {
	if !strings.HasPrefix(pageToken, "hid.") {
		return 0
	}
	return service.HIDToID(pageToken[4:])
}

// IDToPageToken 把 int64 转换成 pageToken
func IDToPageToken(id int64) string {
	if id <= 0 {
		return ""
	}
	return "hid." + service.IDToHID(id)
}

// PageTokenToTime 把 pageToken 转换为 time
func PageTokenToTime(pageToken string) time.Time {
	t := time.Unix(0, 0)
	if pageToken == "" {
		return t
	}
	strs := strings.Split(pageToken, ".")
	if len(strs) != 3 || strs[0] != "unix" {
		return t
	}
	var err error
	var sec int64
	var nsec int64
	if sec, err = strconv.ParseInt(strs[1], 10, 64); err != nil || sec <= 0 {
		return t
	}
	if nsec, err = strconv.ParseInt(strs[2], 10, 64); err != nil || nsec <= 0 {
		return t
	}
	return time.Unix(sec, nsec)
}

// TimeToPageToken 把 time 转换成 pageToken
func TimeToPageToken(t time.Time) string {
	t = t.UTC()
	if t.Unix() <= 0 {
		return ""
	}
	return fmt.Sprintf("unix.%d.%d", t.Unix(), t.Nanosecond())
}
