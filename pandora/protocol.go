package pandora

// PandoraResponseError pandora 报错信息。
type PandoraResponseError struct {
	RequestID string `json:"RequestId"`
	Code      string `json:"Code"`
	Message   string `json:"Message"`
}

func (e *PandoraResponseError) Error() string {
	return e.Message
}

// QueryMode pandora 搜索任务模式，包括极速模式、智能模式、详细模式。
type QueryMode string

const (
	// QueryModeFast 极速模式，不对各字段的字段值做统计分析；同时做字段裁剪，只提取SPL用到的字段
	QueryModeFast QueryMode = "fast"
	// QueryModeSmart 智能模式，不对各字段的字段值做统计分析
	QueryModeSmart QueryMode = "smart"
	// QueryModeDetailed 详细模式，对各字段的字段值做统计分析；同时提取所有字段，启动字段发现
	QueryModeDetailed QueryMode = "detailed"
)

// CreateJobRequest 创建pandora 搜索任务的请求。
type CreateJobRequest struct {
	// Query 搜索的SPL查询语句
	Query string `json:"query"`
	// StartTimeMS 数据开始时间，以毫秒为单位的时间戳。
	StartTimeMS int64 `json:"startTime"`
	// EndTimeMS 数据结束时间，以毫秒为单位的时间戳。
	EndTimeMS int64 `json:"endTime"`
	// Preview 是否提前返回数据
	Preview bool `json:"preview"`
	// CollectSize 限制返回条数。
	CollectSize int `json:"collectSize"`
	// Mode 搜索模式。
	Mode QueryMode `json:"mode"`
}

// CreateJobResponse 创建搜索任务的返回结果，返回搜索任务ID。
type CreateJobResponse struct {
	ID string `json:"id"`
}

// JobProcess 任务进行状态
type JobProcess int

const (
	// JobProcessRunning 0 表示搜索任务进行中
	JobProcessRunning JobProcess = 0
	// JobProcessDone 1 表示搜索任务完成
	JobProcessDone JobProcess = 1
)

// JobStatusResponse 获取搜索任务状态的返回结果。
type JobStatusResponse struct {
	Process    JobProcess `json:"process"`
	DurationMS int64      `json:"duration"`
	EventSize  int64      `json:"eventSize"`
	IsResult   bool       `json:"isResult"`
	IsExport   bool       `json:"isExport"`
	ResultSize int64      `json:"resultSize"`
	ScanSize   int64      `json:"scanSize"`
}

// FieldFlag 字段类型。
type FieldFlag string

const (
	//FieldFlagBucket bucket 代表是分组字段
	FieldFlagBucket FieldFlag = "bucket"
	// FieldFlagMetric metric 代表是统计值
	FieldFlagMetric FieldFlag = "metric"
)

// JobResultsField 搜索任务结果字段属性
type JobResultsField struct {
	Flag        FieldFlag `json:"flag"`
	Name        string    `json:"name"`
	BucketIndex int       `json:"bucketIndex"`
}

// JobResultsResponse 获取搜索任务结果的返回内容。
type JobResultsResponse struct {
	Fields []JobResultsField `json:"fields"`
	Rows   [][]interface{}   `json:"rows"`
}
