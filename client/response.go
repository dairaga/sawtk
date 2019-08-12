package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// response is to save result form http.Reponse.
type response struct {
	code    int
	result  []byte
	lastErr error
}

func (resp *response) String() string {
	return fmt.Sprintf(`{code: %d, result: "%s", last_err: %v}`, resp.code, string(resp.result), resp.lastErr)
}

func (resp *response) fill(httpresp *http.Response) {
	resp.result, resp.lastErr = ioutil.ReadAll(httpresp.Body)
	if resp.lastErr == nil {
		resp.code = httpresp.StatusCode
		httpresp.Body.Close()
	}
}

func (resp *response) handle(expectCode int, data interface{}) error {
	if resp.lastErr != nil {
		return resp.lastErr
	}

	if resp.code == expectCode {
		if err := json.Unmarshal(resp.result, data); err != nil {
			return err
		}
		return nil
	}

	return NewError(resp.code, resp.result)
}
