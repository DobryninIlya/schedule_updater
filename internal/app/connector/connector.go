package connector

import (
	"encoding/json"
	"io"
	"main/internal/app/model"
	"net/http"
)

type Connector struct {
	BaseURL       string
	debug         bool
	ServiceApiKey string
}

func NewConnector(baseURL string, debug bool, serviceApiKey string) *Connector {
	return &Connector{BaseURL: baseURL, debug: debug, ServiceApiKey: serviceApiKey}
}

// GetGroups ...
func (c *Connector) GetGroups() ([]model.Group, error) {
	requestTemplate := c.BaseURL + "/groups" + "?token=" + c.ServiceApiKey
	//variables := fmt.Sprintf(`{"login":"%v"}`, login)
	req, err := http.NewRequest("GET", requestTemplate, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var GroupsAnswer model.GroupAnswer
	var apiResult model.Result
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, &apiResult)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(apiResult.Result, &GroupsAnswer); err != nil {
		return nil, err
	}
	return GroupsAnswer.Data, nil
}

func (c *Connector) GetBuildings() ([]model.Building, error) {
	requestTemplate := c.BaseURL + "/buildings" + "?token=" + c.ServiceApiKey
	req, err := http.NewRequest("GET", requestTemplate, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var BuildingsAnswer model.BuildingsAnswer
	var apiResult model.Result
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, &apiResult)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(apiResult.Result, &BuildingsAnswer); err != nil {
		return nil, err
	}
	return BuildingsAnswer.Data, nil
}

func (c *Connector) GetKafedras() ([]model.Kafedra, error) {
	requestTemplate := c.BaseURL + "/kafedras" + "?token=" + c.ServiceApiKey
	req, err := http.NewRequest("GET", requestTemplate, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var KafedrasAnswer model.KafedraAnswer
	var apiResult model.Result
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, &apiResult)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(apiResult.Result, &KafedrasAnswer); err != nil {
		return nil, err
	}
	return KafedrasAnswer.Data, nil
}

func (c *Connector) GetFaculties() ([]model.Faculty, error) {
	requestTemplate := c.BaseURL + "/faculties" + "?token=" + c.ServiceApiKey
	req, err := http.NewRequest("GET", requestTemplate, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var FacultiesAnswer model.FacultiesAnswer
	var apiResult model.Result
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, &apiResult)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(apiResult.Result, &FacultiesAnswer); err != nil {
		return nil, err
	}
	return FacultiesAnswer.Data, nil
}

func (c *Connector) GetAuditoriums() ([]model.Room, error) {
	requestTemplate := c.BaseURL + "/auditories" + "?token=" + c.ServiceApiKey
	req, err := http.NewRequest("GET", requestTemplate, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var RoomsAnswer model.RoomsAnswer
	var apiResult model.Result
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, &apiResult)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(apiResult.Result, &RoomsAnswer); err != nil {
		return nil, err
	}
	return RoomsAnswer.Data, nil
}

func (c *Connector) GetScheduleByGroup(groupnum string) ([]model.GroupSchedule, error) {
	requestTemplate := c.BaseURL + "/" + groupnum + "?token=" + c.ServiceApiKey
	req, err := http.NewRequest("GET", requestTemplate, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var GroupScheduleAnswer model.GroupScheduleAnswer
	var apiResult model.Result
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, &apiResult)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(apiResult.Result, &GroupScheduleAnswer); err != nil {
		return nil, err
	}
	return GroupScheduleAnswer.Data, nil
}
