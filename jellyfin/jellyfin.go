package jellyfin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Trey2k/helperr/common"
)

type SJellyfin struct {
	URI         string
	PublicURI   string
	AccessToken string
}

func NewConnection() (*SJellyfin, error) {
	jf := &SJellyfin{
		URI:       common.Config.Jellyfin.URI,
		PublicURI: common.Config.Jellyfin.PublicURI,
	}

	authReq := &AuthByNameRequest{
		Username: common.Config.Jellyfin.Username,
		Pw:       common.Config.Jellyfin.Password,
	}

	jsonData, err := json.Marshal(authReq)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/Users/AuthenticateByName", jf.URI), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := jf.doRequest(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jellyfinusername or password is incorrect. %s", http.StatusText(resp.StatusCode))
	}

	response := &AuthenticateByNameResponse{}

	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return nil, err
	}

	jf.AccessToken = response.AccessToken

	return jf, nil
}

func (jf *SJellyfin) doRequest(request *http.Request) (*http.Response, error) {
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("Authorization", `MediaBrowser Client="Helperr", Device="Helperr", DeviceId="None", Version="10.8.5"`)

	client := &http.Client{}

	return client.Do(request)
}

func (jf *SJellyfin) doAuthRequest(request *http.Request) (*http.Response, error) {
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("Authorization", fmt.Sprintf(`MediaBrowser Token="%s", MediaBrowser Client="Helperr", Device="Helperr", DeviceId="None", Version="10.8.5"`, jf.AccessToken))

	client := &http.Client{}

	return client.Do(request)
}

func (jf *SJellyfin) NewUser(name, password string) (*UserInfo, error) {

	userReq := &NewUserRequest{
		Name:     name,
		Password: password,
	}

	jsonData, err := json.Marshal(userReq)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/Users/New", jf.URI), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	resp, err := jf.doAuthRequest(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("none OK status code. %s", http.StatusText(resp.StatusCode))
	}

	toReturn := &UserInfo{}
	err = json.NewDecoder(resp.Body).Decode(toReturn)
	return toReturn, err
}

func (jf *SJellyfin) GetUsers() ([]UserInfo, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/Users", jf.URI), nil)
	if err != nil {
		return nil, err
	}

	resp, err := jf.doAuthRequest(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("none OK status code. %s", http.StatusText(resp.StatusCode))
	}

	toReturn := &[]UserInfo{}
	err = json.NewDecoder(resp.Body).Decode(toReturn)
	return *toReturn, err
}
