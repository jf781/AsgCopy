package operation_test

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"AsgCopy/operation"
)

// Mock data for subscriptions
var mockAvailableSubscriptions = []operation.Subscription{
	{SubscriptionID: "sub1", SubscriptionName: "Subscription 1"},
	{SubscriptionID: "sub2", SubscriptionName: "Subscription 2"},
	{SubscriptionID: "sub3", SubscriptionName: "Subscription 3"},
}

var mockProvideSubscriptions = string("sub1,sub3,unknownSub")

var mockTargetDir = "test-tf-dir"

// Mock data for ASGs
var mockAsgs = map[string][]operation.Asg{
	"sub1": {
		{AsgName: "asg1-sub1", ResourceGroupName: "rg1-1"},
		{AsgName: "asg2-sub1", ResourceGroupName: "rg1-2"},
	},
	"sub2": {
		{AsgName: "asg1-sub2", ResourceGroupName: "rg2-1"},
		{AsgName: "asg2-sub2", ResourceGroupName: "rg2-1"},
		{AsgName: "asg3-sub2", ResourceGroupName: "rg2-2"},
	},
	"sub3": {}, // No ASGs for this subscription
}

// Mock function for SubscriptionLookup
func MockSubscriptionLookup(cred interface{}) []operation.Subscription {
	return mockAvailableSubscriptions
}

// Mock function for AsgLookup
func MockAsgLookup(subscriptionID string, cred interface{}) []operation.Asg {
	return mockAsgs[subscriptionID]
}

func TestSubscriptionLookup(t *testing.T) {
	subList := MockSubscriptionLookup(nil)
	if len(subList) == 0 {
		t.Error("expected non-empty subscription list")
	}

	for _, sub := range subList {
		if sub.SubscriptionID == "" {
			t.Error("expected non-empty subscription ID")
		}
		if sub.SubscriptionName == "" {
			t.Error("expected non-empty subscription name")
		}
	}
}

func TestAsgLookup(t *testing.T) {
	// Replace with a valid subscription ID for testing
	subscriptionID := mockAvailableSubscriptions[0].SubscriptionID
	asgList := MockAsgLookup(subscriptionID, nil)
	if len(asgList) == 0 {
		t.Error("expected non-empty ASG list")
	}

	for _, asg := range asgList {
		if asg.AsgName == "" {
			t.Error("expected non-empty ASG name")
		}
	}
}

func TestPrintSubAsgList(t *testing.T) {
	subName := mockAvailableSubscriptions[0].SubscriptionName
	subAsgList := mockAsgs[mockAvailableSubscriptions[0].SubscriptionID]

	operation.PrintSubAsgList(subName, subAsgList, mockTargetDir)
	formattedSubName := strings.ToLower(strings.ReplaceAll(subName, " ", "-"))

	fileName := formattedSubName + ".tfvars"
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Fatalf("expected file %s to be created", fileName)
	}

	// Clean up the test file
	defer os.Remove(fileName)

	// Additional checks for file content can be added here
}

func TestIntegration(t *testing.T) {
	subList := MockSubscriptionLookup(nil)
	if len(subList) == 0 {
		t.Error("expected non-empty subscription list")
	}

	for _, sub := range subList {
		subAsgList := MockAsgLookup(sub.SubscriptionID, nil)
		if len(subAsgList) == 0 && sub.SubscriptionID != "sub3" {
			t.Errorf("expected non-empty ASG list for subscription %s", sub.SubscriptionID)

			operation.PrintSubAsgList(sub.SubscriptionName, subAsgList, mockTargetDir)
			formattedSubName := strings.ToLower(strings.ReplaceAll(sub.SubscriptionName, " ", "-"))

			fileName := formattedSubName + ".tfvars"
			if _, err := os.Stat(fileName); os.IsNotExist(err) {
				t.Fatalf("expected file %s to be created", fileName)
			}

			// Clean up the test file
			defer os.Remove(fileName)
		}
	}
}

func TestFindMatches(t *testing.T) {
	expectedFoundSubs := []operation.Subscription{
		{SubscriptionID: "sub1", SubscriptionName: "Subscription 1"},
		{SubscriptionID: "sub3", SubscriptionName: "Subscription 3"},
	}

	expectedMissingSubs := []string{"unknownSub"}

	foundSubs, missingSubs := operation.FindMatches(mockAvailableSubscriptions, mockProvideSubscriptions)

	if !reflect.DeepEqual(foundSubs, expectedFoundSubs) {
		t.Errorf("expected foundSubs to be %v, got %v", expectedFoundSubs, foundSubs)
	}

	if !reflect.DeepEqual(missingSubs, expectedMissingSubs) {
		t.Errorf("expected missingSubs to be %v, got %v", expectedMissingSubs, missingSubs)
	}
}

func TestEnsureDir(t *testing.T) {

	// Test creating a new directory
	err := operation.EnsureDir(mockTargetDir)
	if err != nil {
		t.Fatalf("EnsureDir failed to create directory: %v", err)
	}

	info, err := os.Stat(mockTargetDir)
	if err != nil {
		t.Fatalf("Failed to stat the directory: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("Path exists but is not a directory: %s", mockTargetDir)
	}

	// Test running EnsureDir on an existing directory
	err = operation.EnsureDir(mockTargetDir)
	if err != nil {
		t.Fatalf("EnsureDir failed on an existing directory: %v", err)
	}

	os.RemoveAll(mockTargetDir)

}
