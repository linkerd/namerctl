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
		List() ([]string, error)
		Get(name string) (*VersionedDtab, error)
		Create(name string, dtab Dtab) (Version, error)
		Update(name string, dtab Dtab, version Version) (Version, error)
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

func (ctl *httpController) List() ([]string, error) {
	req, err := ctl.dtabRequest("GET", "", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	rsp, err := ctl.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer drainAndClose(rsp)

	switch rsp.StatusCode {
	case http.StatusOK:
		var names []string
		if err := json.NewDecoder(rsp.Body).Decode(&names); err != nil {
			return nil, err
		}
		return names, nil

	default:
		return nil, fmt.Errorf("unexpected response: %s", rsp.Status)
	}
}

func (ctl *httpController) Get(name string) (*VersionedDtab, error) {
	req, err := ctl.dtabRequest("GET", name, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/dtab")

	rsp, err := ctl.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer drainAndClose(rsp)

	switch rsp.StatusCode {
	case http.StatusOK:
		v := Version(rsp.Header.Get("ETag"))
		bytes, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}
		dtab, err := ParseDtab(string(bytes))
		if err != nil {
			return nil, err
		}
		return &VersionedDtab{v, dtab}, nil

	default:
		return nil, fmt.Errorf("unexpected response: %s", rsp.Status)
	}
}

func (ctl *httpController) Create(name string, dtab Dtab) (Version, error) {
	req, err := ctl.dtabRequest("POST", name, strings.NewReader(dtab.String()))
	if err != nil {
		return Version(""), err
	}
	req.Header.Set("Content-Type", "application/dtab")

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

func (ctl *httpController) Update(name string, dtab Dtab, version Version) (Version, error) {
	req, err := ctl.dtabRequest("PUT", name, strings.NewReader(dtab.String()))
	if err != nil {
		return Version(""), err
	}
	req.Header.Set("Content-Type", "application/dtab")
	req.Header.Set("If-Match", string(version))

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
