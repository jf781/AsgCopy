package operation

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

type Subscription struct {
	SubscriptionID   string
	SubscriptionName string
}

type Asg struct {
	AsgName           string
	ResourceGroupName string
}

func SubscriptionLookup(cred azcore.TokenCredential) []Subscription {

	ctx := context.Background()

	var subList []Subscription

	client, err := armsubscriptions.NewClient(cred, nil)
	if err != nil {
		log.Fatalf("Failed to create the subscription client: %v", err)
	}

	pager := client.NewListPager(nil)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to get next page of subscriptions: %v", err)
		}

		for _, sub := range page.Value {
			subList = append(subList, Subscription{SubscriptionID: *sub.SubscriptionID, SubscriptionName: *sub.DisplayName})
		}

	}
	return subList
}

func AsgLookup(subscription Subscription, cred azcore.TokenCredential) []Asg {

	ctx := context.Background()

	asgClient, err := armnetwork.NewApplicationSecurityGroupsClient(subscription.SubscriptionID, cred, nil)
	if err != nil {
		log.Fatalf("Failed to create ASG client: %v", err)
	}

	pager := asgClient.NewListAllPager(nil)

	var asgList []Asg
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("Failed to get next page of ASGs: %v", err)
		}

		for _, asg := range page.Value {
			resourceGroupName := strings.Split(*asg.ID, "/")[4]
			asgList = append(asgList, Asg{AsgName: *asg.Name, ResourceGroupName: resourceGroupName})
		}
	}
	return asgList
}

func PrintSubAsgList(subName string, subAsgList []Asg, targetDir string) error {

	if len(subAsgList) != 0 {

		os.Chdir(targetDir)
		formattedSubName := strings.ToLower(strings.ReplaceAll(subName, " ", "-"))

		file, err := os.Create(formattedSubName + ".tf")
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
		defer file.Close()

		_, err = file.WriteString("asgs = [\n")
		if err != nil {
			log.Fatalf("failed to write to file: %v", err)
		}

		for _, asg := range subAsgList {
			// Formatting the block of Terraform code
			block := `  { 
		asgName           = "` + asg.AsgName + `"
		resourceGroupName = "` + asg.ResourceGroupName + `"
	},` + "\n"
			_, err = file.WriteString(block)
			if err != nil {
				log.Fatalf("failed to write to file: %v", err)
				return err
			}
		}

		_, err = file.WriteString("]\n")
		if err != nil {
			log.Fatalf("failed to write to file: %v", err)

		}
	}

	return nil
}

func FindMatches(availableSubs []Subscription, providedSubs string) (foundSubs []Subscription, missingSubs []string) {
	set := make(map[string]struct{}, len(availableSubs))
	for _, v := range availableSubs {
		set[v.SubscriptionID] = struct{}{}
	}

	providedSubsList := strings.Split(providedSubs, ",")

	var matchingSubs []string
	for _, v := range providedSubsList {
		if _, exists := set[v]; exists {
			matchingSubs = append(matchingSubs, v)
		} else {
			missingSubs = append(missingSubs, v)
		}
	}

	for _, v := range availableSubs {
		for _, m := range matchingSubs {
			if v.SubscriptionID == m {
				foundSubs = append(foundSubs, v)
			}
		}
	}

	return
}

func EnsureDir(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755) // Creates all necessary parent directories
		if err != nil {
			log.Fatalf("failed to create directory: %v", err)
			return err
		}
	} else if err != nil {
		log.Fatalf("error checking directory: %v", err)
		return err
	} else if !info.IsDir() {
		// Path exists but is not a directory
		log.Fatalf("path exists but is not a directory: %s", path)
		return err
	}
	return nil
}
