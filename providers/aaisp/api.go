package aaisp

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	"strings"
)

func (provider *Provider) sendSMS(form url.Values) error {
	req, err := http.NewRequest("POST", "https://sms.aa.net.uk/sms.cgi", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	respBytes, err := ioutil.ReadAll(httpResp.Body)
	httpResp.Body.Close()
	if err != nil {
		return fmt.Errorf("Error reading response from AAISP: %s", err)
	}

	if !(httpResp.StatusCode >= 200 && httpResp.StatusCode <= 299) {
		return fmt.Errorf("HTTP error from AAISP: %s", httpResp.Status)
	}

	resp := string(respBytes)
	if !strings.HasPrefix(resp, "OK") {
		return errors.New(resp)
	}

	return nil
}
