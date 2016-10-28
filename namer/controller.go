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
		Version Version `json:"version"`
		Dtab    Dtab    `json:"dtab"`
	}

	Controller interface {
		List() ([]string, error)
		Get(name string) (*VersionedDtab, error)
		Create(name string, dtabstr string) (Version, error)
		Delete(name string) error
		Update(name string, dtabstr string, version Version) (Version, error)
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
	req.Header.Set("Accept", "application/json")

	rsp, err := ctl.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer drainAndClose(rsp)

	switch rsp.StatusCode {
	case http.StatusOK:
		dtab := VersionedDtab{
			Version: Version(rsp.Header.Get("ETag")),
		}
		if err := json.NewDecoder(rsp.Body).Decode(&dtab.Dtab); err != nil {
			return nil, err
		}
		return &dtab, nil

	default:
		return nil, fmt.Errorf("unexpected response: %s", rsp.Status)
	}
}

func isJson(str string) bool {
	return len(str) > 0 && (str[0:1] == "{" || str[0:1] == "[")
}

func (ctl *httpController) Create(name, dtabstr string) (Version, error) {
	emptyVersion := Version("")
	var req *http.Request
	if isJson(dtabstr) {
		var vdtab VersionedDtab
		if err := json.Unmarshal([]byte(dtabstr), &vdtab); err != nil {
			return emptyVersion, err
		}
		dtab, err := json.Marshal(vdtab.Dtab)
		if err != nil {
			return emptyVersion, err
		}
		req, err = ctl.dtabRequest("POST", name, strings.NewReader(string(dtab)))
		if err != nil {
			return emptyVersion, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err := ctl.dtabRequest("POST", name, strings.NewReader(dtabstr))
		if err != nil {
			return emptyVersion, err
		}
		req.Header.Set("Content-Type", "application/dtab")
	}

	rsp, err := ctl.client.Do(req)
	if err != nil {
		return emptyVersion, err
	}
	defer drainAndClose(rsp)

	switch rsp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		v := Version(rsp.Header.Get("ETag"))
		return v, nil

	default:
		return emptyVersion, fmt.Errorf("unexpected response: %s", rsp.Status)
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

func (ctl *httpController) Update(name, dtabstr string, version Version) (Version, error) {
	useJson := isJson(dtabstr)
	if useJson {
		var vdtab VersionedDtab
		if err := json.Unmarshal([]byte(dtabstr), &vdtab); err != nil {
			return Version(""), err
		}

		dtab, err := json.Marshal(vdtab.Dtab)
		if err != nil {
			return Version(""), err
		}

		dtabstr = string(dtab)
		if vdtab.Version != "" {
			version = vdtab.Version
		}
	}

	req, err := ctl.dtabRequest("PUT", name, strings.NewReader(dtabstr))
	if err != nil {
		return Version(""), err
	}
	if version != "" {
		req.Header.Set("If-Match", string(version))
	}

	if useJson {
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

func drainAndClose(resp *http.Response) {
	if resp != nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
}
