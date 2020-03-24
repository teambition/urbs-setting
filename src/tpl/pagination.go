package tpl

import (
	"net/url"
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
func (pg *Pagination) TokenToTime(defaultTime ...time.Time) time.Time {
	return PageTokenToTime(pg.PageToken, defaultTime...)
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
func PageTokenToTime(pageToken string, defaultTime ...time.Time) time.Time {
	t := time.Unix(0, 0)
	if len(defaultTime) > 0 {
		t = defaultTime[0]
	}
	if pageToken == "" {
		return t
	}

	t2, err := time.Parse(time.RFC3339, pageToken)
	if err != nil {
		return t
	}
	return t2
}

// TimeToPageToken 把 time 转换成 pageToken
func TimeToPageToken(t time.Time) string {
	t = t.UTC()
	if t.Unix() <= 0 {
		return ""
	}
	return url.QueryEscape(t.Format(time.RFC3339))
}
