package tpl

// ResponseType 定义了标准的 List 接口返回数据模型
type ResponseType struct {
	Error         string      `json:"error,omitempty"`
	Message       string      `json:"message,omitempty"`
	NextPageToken string      `json:"nextPageToken,omitempty"`
	TotalSize     uint64      `json:"totalSize,omitempty"`
	Result        interface{} `json:"result,omitempty"`
}
