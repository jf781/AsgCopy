# ASG Lookup

This project will create a utility to list all Application Security Groups (ASGs) in defined Azure subscriptions and create a Terraform file that will create the ASGs in a different subscription. 

The application will create a Terraform file with a list of ASGs that are present in the source subscription.  The Terraform files will be named after the subscription and will contain the ASGs in the following format.
  
  ```hcl
  prod_asgs = [
    { 
      name           = "asg1"
      resource_group_name = "rg1"
    },
  ]
  ```

It requires a list of subscription IDs to fetch the ASGs from. It will confirm that the subscription IDs are valid and that the user has access to them. If the subscription IDs are not valid, the application will fetch the ASGs from the subscriptions that are valid.  It will output a list of the subscriptions that were not found, or the user does not have access to.

### Flags
- `--subscriptionsIds`: Comma-separated list of subscription IDs to fetch ASGs from.
- `--targetDir`: The directory where the Terraform files will be created. Default is `./tf-files`.

## Usage

To run the application, use the following command:
```sh
go run [main.go](http://_vscodecontentref_/3) --subscriptionsIds=<comma-separated-subscription-ids>  --targetDir=<target-directory>
```
## Project Structure
```
.
├── main.go
└── operation/
  ├── asgLookup.go
  └── asgLookup_test.go
```

- `main.go`: The entry point of the application.
- `operation/asgLookup.go`: Contains the core logic for subscription and ASG lookup.
- `operation/asgLookup_test.go`: Contains unit tests for the functions in `asgLookup.go`.

## Installation

1. Clone the repository:
    ```sh
    git clone <repository-url>
    cd <repository-directory>
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

## Functions
### SubscriptionLookup
Fetches the list of subscriptions using Azure SDK.

### AsgLookup
Fetches the list of ASGs for a given subscription using Azure SDK.

### FindMatches
Finds matching and missing subscriptions from the provided list.

### PrintSubAsgList
Prints the list of ASGs for a subscription to a Terraform file.

### ensureDir
Validates if the directory exists and creates it if it doesn't.

## Testing
To run the test against the functions, use the following command:

```sh
go test -v ./operation
```
