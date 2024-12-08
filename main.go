package main

import (
	"log"
	"os"
	"fmt"


	"AsgCopy/operation"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:    "asgLookup",
		Version: "Development",
		Usage:   "This too lists all ASGs in defined subscriptions and will create a Terraform file with the ASGs",
		Action: run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "subscriptionsIds",
				Usage:    "Provide a comma-separated list of subscription IDs to query for ASGs.",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "targetDir",
				Usage:    "Directory to output the Terraform files.",
				Required: false,
				Value:    "tf-files",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}


func run(c *cli.Context) error {
	var subAsgList []operation.Asg
	
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	
	availableSubscriptions := operation.SubscriptionLookup(cred)

	matched, notFound :=operation.FindMatches(availableSubscriptions, c.String("subscriptionsIds"))

	err = operation.EnsureDir(c.String("targetDir"))
	if err != nil {
		log.Fatalf("failed to access target directory: %v", err)
	}

	os.Chdir("tf-files")


	for _, sub := range matched {

		subAsgList = operation.AsgLookup(sub, cred)
		operation.PrintSubAsgList(sub.SubscriptionName, subAsgList, c.String("targetDir"))
	}

	if len(notFound) > 0 {
		fmt.Printf("The following subscriptions were not found: \n %v", notFound)
	}
	return nil
}