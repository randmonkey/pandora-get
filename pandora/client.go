package pandora

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	server string
	token  string

	client *http.Client
}

type pandoraTokenTransport struct {
	token string
	tr    http.RoundTripper
}

func (tr *pandoraTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", tr.token)
	return tr.tr.RoundTrip(req)
}

func NewClient(server string, token string) *Client {
	c := &Client{
		server: strings.TrimSuffix(server, "/"),
		client: &http.Client{
			Transport: &pandoraTokenTransport{token: token, tr: http.DefaultTransport},
		},
	}
	return c
}

func parseResponseError(respBody []byte) error {
	respErr := &PandoraResponseError{}
	parseErr := json.Unmarshal(respBody, &respErr)
	if parseErr != nil {
		return parseErr
	}
	return respErr
}

func (c *Client) CreateJob(args *CreateJobRequest) (*CreateJobResponse, error) {
	buf, _ := json.Marshal(args)
	resp, err := c.client.Post(c.server+"/api/v1/jobs", "application/json", bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	if resp.StatusCode/100 != 2 {
		return nil, parseResponseError(buf)
	}
	ret := &CreateJobResponse{}
	err = json.Unmarshal(buf, ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func (c *Client) GetJobStatus(jobID string) (*JobStatusResponse, error) {
	resp, err := c.client.Get(c.server + "/api/v1/jobs/" + jobID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	if resp.StatusCode/100 != 2 {
		return nil, parseResponseError(buf)
	}
	ret := &JobStatusResponse{}
	err = json.Unmarshal(buf, ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func (c *Client) GetJobResults(jobID string) (*JobResultsResponse, error) {
	resp, err := c.client.Get(c.server + "/api/v1/jobs/" + jobID + "/results")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	if resp.StatusCode/100 != 2 {
		return nil, parseResponseError(buf)
	}
	ret := &JobResultsResponse{}
	err = json.Unmarshal(buf, ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

// ToKVSlice 将pandora的返回结果转换成 key-value 数组切片。
func (r *JobResultsResponse) ToKVSlice() []map[string]interface{} {
	if r == nil {
		return nil
	}
	if len(r.Fields) == 0 || len(r.Rows) == 0 {
		return nil
	}

	ret := make([]map[string]interface{}, 0)
	for _, row := range r.Rows {
		m := make(map[string]interface{})
		for i, field := range r.Fields {
			m[field.Name] = row[i]
		}
		ret = append(ret, m)
	}
	return ret
}

func (c *Client) GetQueryResult(spl string, startTime, endTime time.Time, limit int, timeout time.Duration) ([]map[string]interface{}, error) {
	jobReq := &CreateJobRequest{
		Query:       spl,
		StartTimeMS: startTime.UnixNano() / (1000 * 1000),
		EndTimeMS:   endTime.UnixNano() / (1000 * 1000),
		Preview:     false,
		CollectSize: limit,
		Mode:        QueryModeFast,
	}

	jobRes, err := c.CreateJob(jobReq)
	if err != nil {
		log.Printf("failed to create job, error %v", err)
		return nil, err
	}

	jobID := jobRes.ID

	// 开启任务轮询计时器与任务超时计时器。
	if timeout == 0 {
		timeout = 60 * time.Second
	}
	pollInterval := time.Second

	pollTicker := time.NewTicker(pollInterval)
	expireTimer := time.NewTimer(timeout)
	defer pollTicker.Stop()
	defer expireTimer.Stop()
	var jobResults *JobResultsResponse
	for {
		select {
		case <-pollTicker.C:
			status, err := c.GetJobStatus(jobID)
			if err != nil {
				return nil, err
			}
			if status.Process == JobProcessDone {
				results, err := c.GetJobResults(jobID)
				if err != nil {
					return nil, err
				}
				jobResults = results
				return jobResults.ToKVSlice(), nil
			}
		case <-expireTimer.C:
			return nil, fmt.Errorf("job timout")
		}
	}

}
