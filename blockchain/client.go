package blockchain

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

type TranItem struct {
	Tran		string	`json:"transaction"`
	TimeStamp	int64	`json:"timestamp"`
}

type TranItems []TranItem
func (g TranItems) Len() int {
	return len(g)
}
func (g TranItems) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}
func (g TranItems) Less(i, j int) bool {
	return g[i].TimeStamp < g[j].TimeStamp
}

type getTransResponse struct {
	Limit	uint64		`json:"limit"`
	Offset	uint64		`json:"offset"`
	Total	uint64		`json:"total"`
	Trans	TranItems	`json:"transactions"`
}

type paramsMap map[string]string
func newParamsMap() paramsMap {
	return make(paramsMap)
}
func (p paramsMap) ToString() string {
	params := url.Values{}

	for key, value := range p {
		params.Add(key, value)
	}

	return params.Encode()
}

type Trans struct {
	Data	TranItems
	ReadAll	bool
}

type Client struct {
	baseURL				*url.URL
	userAgent			string

	httpClient			*http.Client

	unSendTrans			[]*Tran
	muSendTrans			sync.Mutex

	getTranExecute		*int32
	lastTranTimestamp	*int64
	oldLastTran			string
	muUpdateOldLastTran	sync.RWMutex
}

func NewClient(rawURL string) (*Client, error) {
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	var getTranExecute int32 = 0
	var lastTranTimestamp int64 = 0

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{},
		lastTranTimestamp: &lastTranTimestamp,
		getTranExecute: &getTranExecute,
		oldLastTran: "",
	}, nil
}

func (c *Client) SendTrans(trans []*Tran) error {
	if trans == nil || len(trans) == 0 {
		return nil
	}
	fullTrans := c.createTrans(trans)
	fullTranStrings := make([]string, 0, len(fullTrans))

	for i, n := 0, len(fullTrans); i < n; i ++ {
		s := fullTrans[i].ToString()
		if s != "" {
			fullTranStrings = append(fullTranStrings, s)
		}
	}

	req, err := c.newRequest(
		"POST",
		"/AddTransactions",
		nil,
		fullTranStrings,
	)
	if err != nil {
		c.saveTrans(fullTrans)
		return err
	}

	_, _, err = c.do(req)
	if err != nil {
		c.saveTrans(fullTrans)
		return err
	}

	return nil
}

func (c *Client) GetTrans() (*Trans, error) {
	if !atomic.CompareAndSwapInt32(c.getTranExecute, 0, 1) {
		return nil, errors.New("get tran query executed yet")
	} else {
		defer func() {
			atomic.StoreInt32(c.getTranExecute, 0)
		}()
	}

	params := newParamsMap()
	params["from"] = strconv.FormatInt(atomic.LoadInt64(c.lastTranTimestamp), 10)
	params["offset"] = strconv.FormatInt(0, 10)
	params["limit"] = strconv.FormatInt(1000, 10)
	params["order"] = "older"

	req, err := c.newRequest(
		"GET",
		"/GetTransactions",
		params,
		nil,
	)
	if err != nil {
		return nil, err
	}

	_, rawTrans, err := c.do(req)
	if err != nil {
		return nil, err
	}
	var transResp getTransResponse
	err = json.Unmarshal(rawTrans, &transResp)
	if err != nil {
		return nil, err
	}

	trans := transResp.Trans

	if trans == nil || len(trans) == 0 {
		return &Trans{
			Data: nil,
			ReadAll: true,
		}, nil
	}

	sort.Sort(trans)
	oldTimestamp := atomic.LoadInt64(c.lastTranTimestamp)
	oldTran := c.getLastTran()

	for i, tran := range trans {
		if tran.TimeStamp == oldTimestamp {
			if tran.Tran == oldTran {
				if i + 1 == len(trans) {
					return &Trans{
						Data: nil,
						ReadAll: true,
					}, nil
				}

				trans = trans[(i+1):]
				break
			}
		} else {
			break
		}
	}

	c.updateOldLastTran(&trans[len(trans) - 1])

	return &Trans{
		Data: trans,
		ReadAll: transResp.Total <= transResp.Offset + transResp.Limit,
	}, nil
}

func (c *Client) getLastTran() string {
	c.muUpdateOldLastTran.RLock()
	defer c.muUpdateOldLastTran.RUnlock()

	return c.oldLastTran
}


func (c *Client) updateOldLastTran(tran *TranItem) {
	c.muUpdateOldLastTran.Lock()
	defer c.muUpdateOldLastTran.Unlock()

	atomic.StoreInt64(c.lastTranTimestamp, tran.TimeStamp)
	c.oldLastTran = tran.Tran
}

func (c *Client) createTrans(trans []*Tran) []*Tran {
	c.muSendTrans.Lock()
	defer c.muSendTrans.Unlock()

	var fullTrans []*Tran
	if c.unSendTrans != nil {
		fullTrans = append(trans, c.unSendTrans...)
		c.unSendTrans = nil
	} else {
		fullTrans = trans
	}

	return fullTrans
}

func (c *Client) saveTrans(trans []*Tran) {
	c.muSendTrans.Lock()
	defer c.muSendTrans.Unlock()

	if c.unSendTrans == nil {
		c.unSendTrans = trans
	} else {
		c.unSendTrans = append(c.unSendTrans, trans...)
	}
}

func (c *Client) newRequest(
	method string,
	path string,
	params paramsMap,
	body interface{},
) (*http.Request, error) {
	rel := &url.URL{
		Path: path,
	}
	if params != nil {
		rel.RawQuery = params.ToString()
	}
	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json, text/*")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	return req, nil
}

func (c *Client) do(
	req *http.Request,
) (*http.Response, []byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return resp, body, err
}
