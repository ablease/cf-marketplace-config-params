package cf

import (
	"encoding/json"
	"fmt"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

type Client struct {
	cliConnection plugin.CliConnection
}

func NewClient(cliConnection plugin.CliConnection) Client {
	return Client{
		cliConnection,
	}
}

type Results struct {
	Resources []Result
}

type Result struct {
	Metadata struct {
		Guid string
	}
	Entity struct {
		Label string
		Name  string
		Extra string
	}
}

func (c *Client) GetServiceGuid(name string) (string, error) {
	url := "/v2/services?results-per-page=100"

	var results Results
	err := c.getResult(url, &results)

	if err != nil {
		return "", err
	}

	for _, service := range results.Resources {
		if service.Entity.Label == name {
			return service.Metadata.Guid, nil
		}
	}

	return "", fmt.Errorf("Could not find service: %s", name)
}

func (c *Client) GetPlanGuid(name, serviceGuid string) (string, error) {
	url := fmt.Sprintf("/v2/service_plans?results-per-page=100&q=service_guid:%v", serviceGuid)

	var results Results
	err := c.getResult(url, &results)

	if err != nil {
		return "", err
	}

	for _, servicePlan := range results.Resources {
		if servicePlan.Entity.Name == name {
			return servicePlan.Metadata.Guid, nil
		}
	}

	return "", fmt.Errorf("Could not find plan: %s", name)
}

func (c *Client) GetSchema(planGuid string) (string, error) {
	url := fmt.Sprintf("/v2/service_plans/%s", planGuid)

	var result Result
	err := c.getResult(url, &result)

	if err != nil {
		return "", err
	}

	return result.Entity.Extra, nil
}

func (c *Client) getResult(url string, unmarshal interface{}) error {
	body, err := c.cliConnection.CliCommandWithoutTerminalOutput("curl", url)

	raw := strings.Join(body, " ")

	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(raw), &unmarshal)

	if err != nil {
		return err
	}

	return nil
}
