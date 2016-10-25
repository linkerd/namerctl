package namer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// XXX later we should add support for parsing dtabs

type (
	Version       string
	VersionedDtab struct {
		Version Version
		Dtab    Dtab
	}

	Controller interface {
		List() ([]string, string, error)
		Get(name string, requestJson bool) (*VersionedDtab, string, error)
		Create(name string, dtabstr string, isJson bool) (Version, error)
		Delete(name string) error
		Update(name string, dtabstr string, isJson bool, version Version) (Version, error)
	}

	httpController struct {
		baseURL *url.URL
		client  *http.Client
	}
)

func NewHttpController(baseURL *url.URL, client *http.Client) Controller {
	return &httpController{baseURL, client}
}

func (ctl *httpController) dtabRequest(method, name string, data io.Reader) (*http.Request, error) {
	u := *ctl.baseURL
	if name == "" {
		u.Path = "/api/1/dtabs"
	} else {
		u.Path = fmt.Sprintf("/api/1/dtabs/%s", name)
	}
	return http.NewRequest(method, u.String(), data)
}

func (ctl *httpController) List() ([]string, string, error) {
	req, err := ctl.dtabRequest("GET", "", nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Accept", "application/json")

	rsp, err := ctl.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer drainAndClose(rsp)

	switch rsp.StatusCode {
	case http.StatusOK:
		var names []string
		bytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return nil, "", err
		}

		if err = json.Unmarshal(bytes, &names); err != nil {
			return nil, "", err
		}

		return names, string(bytes), nil

	default:
		return nil, "", fmt.Errorf("unexpected response: %s", rsp.Status)
	}
}

func (ctl *httpController) Get(name string, requestJson bool) (*VersionedDtab, string, error) {
	req, err := ctl.dtabRequest("GET", name, nil)
	if err != nil {
		return nil, "", err
	}

	if requestJson {
		req.Header.Set("Accept", "application/json")
	} else {
		req.Header.Set("Accept", "application/dtab")
	}

	rsp, err := ctl.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer drainAndClose(rsp)

	switch rsp.StatusCode {
	case http.StatusOK:
		v := Version(rsp.Header.Get("ETag"))
		bytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return nil, "", err
		}

		if requestJson {
			return &VersionedDtab{v, nil}, string(bytes), nil
		} else {
			dtab, err := parseDtab(string(bytes))
			if err != nil {
				return nil, "", err
			}
			return &VersionedDtab{v, dtab}, "", nil
		}

	default:
		return nil, "", fmt.Errorf("unexpected response: %s", rsp.Status)
	}
}

func (ctl *httpController) Create(name, dtabstr string, isJson bool) (Version, error) {
	req, err := ctl.dtabRequest("POST", name, strings.NewReader(dtabstr))
	if err != nil {
		return Version(""), err
	}

	if isJson {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/dtab")
	}

	rsp, err := ctl.client.Do(req)
	if err != nil {
		return Version(""), err
	}
	defer drainAndClose(rsp)

	switch rsp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		v := Version(rsp.Header.Get("ETag"))
		return v, nil

	default:
		return Version(""), fmt.Errorf("unexpected response: %s", rsp.Status)
	}
}

func (ctl *httpController) Delete(name string) error {

	req, err := ctl.dtabRequest("DELETE", name, nil)
	if err != nil {
		return err
	}
	rsp, err := ctl.client.Do(req)
	if err != nil {
		return err
	}
	defer drainAndClose(rsp)

	switch rsp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		return nil

	default:
		return fmt.Errorf("unexpected response: %s", rsp.Status)
	}
}

func (ctl *httpController) Update(name, dtabstr string, isJson bool, version Version) (Version, error) {
	req, err := ctl.dtabRequest("PUT", name, strings.NewReader(dtabstr))
	if err != nil {
		return Version(""), err
	}
	if isJson {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/dtab")
	}
	if version != "" {
		req.Header.Set("If-Match", string(version))
	}

	rsp, err := ctl.client.Do(req)
	if err != nil {
		return Version(""), err
	}
	defer drainAndClose(rsp)

	switch rsp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		v := Version(rsp.Header.Get("ETag"))
		return v, nil

	default:
		return Version(""), fmt.Errorf("unexpected response: %s", rsp.Status)
	}
}

func drainAndClose(resp *http.Response) {
	if resp != nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
}
