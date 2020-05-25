package tpl

import (
	"strings"
	"time"

	"github.com/teambition/gear"
	"github.com/teambition/urbs-setting/src/service"
)

// Search 搜索
type Search struct {
	Q string `json:"q" query:"q"`
}

// Validate escape and build MySQL LIKE pattern
func (s *Search) Validate() error {
	if s.Q != "" {
		if len(s.Q) <= 2 {
			return gear.ErrBadRequest.WithMsgf("too small query: %s", s.Q)
		}
		s.Q = strings.ReplaceAll(s.Q, `\`, "-")
		s.Q = strings.ReplaceAll(s.Q, "%", `\%`)
		s.Q = strings.ReplaceAll(s.Q, "_", `\_`)
	}
	if s.Q != "" {
		s.Q = s.Q + "%" // %q% 在大数据表（如user表）下开销太大
	}
	return nil
}

// Pagination 分页
type Pagination struct {
	Search
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
		return gear.ErrBadRequest.WithMsgf("pageSize %v should not great than 1000", pg.PageSize)
	}

	if pg.PageSize <= 0 {
		pg.PageSize = 10
	}

	if err := pg.Search.Validate(); err != nil {
		return err
	}

	return nil
}

// TokenToID 把 pageToken 转换为 int64
func (pg *Pagination) TokenToID() int64 {
	return PageTokenToID(pg.PageToken)
}

// TokenToTimestamp 把 pageToken 转换为 timestamp (ms)
func (pg *Pagination) TokenToTimestamp(defaultTime ...time.Time) int64 {
	return PageTokenToTimestamp(pg.PageToken, defaultTime...)
}

// PageTokenToID 把 pageToken 转换为 int64
func PageTokenToID(pageToken string) int64 {
	if !strings.HasPrefix(pageToken, "h.") {
		return 9223372036854775807
	}
	return service.HIDToID(pageToken[2:])
}

// IDToPageToken 把 int64 转换成 pageToken
func IDToPageToken(id int64) string {
	if id <= 0 {
		return ""
	}
	return "h." + service.IDToHID(id)
}

// PageTokenToTimestamp 把 pageToken 转换为 timestamp
func PageTokenToTimestamp(pageToken string, defaultTime ...time.Time) int64 {
	t := time.Unix(0, 0)
	if len(defaultTime) > 0 {
		t = defaultTime[0]
	}
	if !strings.HasPrefix(pageToken, "t.") {
		return t.Unix()*1000 + int64(t.UTC().Nanosecond()/1000000)
	}

	return service.HIDToID(pageToken[2:])
}

// TimeToPageToken 把 time 转换成 pageToken
func TimeToPageToken(t time.Time) string {
	s := t.Unix()*1000 + int64(t.UTC().Nanosecond()/1000000)
	if s <= 0 {
		return ""
	}
	return "t." + service.IDToHID(s)
}
