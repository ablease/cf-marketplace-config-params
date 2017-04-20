package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"code.cloudfoundry.org/cli/plugin"

	"github.com/fatih/color"
	"github.com/pivotal-cf/marketplacev2/cf"
)

type marketplaceplugin struct{}

type SchemaResult struct {
	Schemas struct {
		ServiceInstances struct {
			Create struct {
				Parameters map[string]interface{}
			}
		} `json:"service_instances"`
	}
}

func (m *marketplaceplugin) Run(cliConnection plugin.CliConnection, args []string) {
	if len(args) > 0 && args[0] == "CLI-MESSAGE-UNINSTALL" {
		return
	}

	fmtGreen := color.New(color.FgGreen, color.Bold)

	args = args[1:]

	f := flag.NewFlagSet("marketplace-v2", flag.ExitOnError)

	var service string
	var plan string

	f.StringVar(&service, "s", "", "")
	f.StringVar(&plan, "p", "", "")

	err := f.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	if service == "" && plan == "" {
		_, err := cliConnection.CliCommand("m")
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if service == "" || plan == "" {
		log.Fatal("Must supply both a service and a plan name")
	}

	fmt.Printf("Getting configuration parameter schemas for service %s and plan %s\n", service, plan)

	schemaResult, err := m.GetPlanSchema(service, plan, cliConnection)
	if err != nil {
		log.Fatal(err)
	}

	fmtGreen.Print("OK\n\n")

	schemaParams := schemaResult.Schemas.ServiceInstances.Create.Parameters

	if schemaParams == nil {
		fmt.Printf("Plan %s does not support configuration parameter schemas\n", plan)
		os.Exit(0)
	}

	out, err := json.MarshalIndent(schemaParams, "", "  ")

	fmt.Printf("Create Service Configuration Parameters:\n\n")
	fmt.Println(string(out))
}

func (m *marketplaceplugin) GetPlanSchema(service, plan string, cliConnection plugin.CliConnection) (SchemaResult, error) {
	var schemaResult SchemaResult

	c := cf.NewClient(cliConnection)

	serviceGuid, err := c.GetServiceGuid(service)

	if err != nil {
		return schemaResult, err
	}

	planGuid, err := c.GetPlanGuid(plan, serviceGuid)
	if err != nil {
		return schemaResult, err
	}

	schema, err := c.GetSchema(planGuid)
	if err != nil {
		return schemaResult, err
	}

	err = json.Unmarshal([]byte(schema), &schemaResult)
	if err != nil {
		return schemaResult, err
	}

	return schemaResult, nil
}

func (m *marketplaceplugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "marketplace-v2",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "marketplace-v2",
				HelpText: "List available offerings in the marketplace",

				UsageDetails: plugin.Usage{
					Usage: "cf marketplace-v2 [-s SERVICE [ -p PLAN ]]",
					Options: map[string]string{
						"s": "Show plan details for a particular service offering",
						"p": "Show any configuration parameter schemas for a particular plan and service offering",
					},
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(marketplaceplugin))
}
