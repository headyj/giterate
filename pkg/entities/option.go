package entities

import (
	"log"
	"strconv"

	"github.com/go-git/go-git/plumbing"

	"github.com/go-git/go-git"
)

type Option struct {
	Name  string `json:"Name" yaml:"Name"`
	Value string `json:"Value" yaml:"Value"`
}

func ProcessCloneOptions(options *[]Option, cloneOptions *git.CloneOptions) {
	var err error
	for _, option := range *options {
		switch option.Name {
		case "SingleBranch":
			cloneOptions.SingleBranch, err = strconv.ParseBool(option.Value)
		case "NoCheckout":
			cloneOptions.NoCheckout, err = strconv.ParseBool(option.Value)
		case "RemoteName":
			cloneOptions.RemoteName = option.Value
		case "Branch", "ReferenceName":
			cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(option.Value)
		case "Depth":
			i, err := strconv.ParseInt(option.Value, 10, 32)
			if err != nil {
				log.Fatalf("Cannot convert Depth value to int: %s\n", err)
			}
			cloneOptions.Depth = int(i)
		default:
			log.Fatalf("%s is not a valid CloneOption\n", option.Name)
		}
		if err != nil {
			log.Fatalf("Cannot get %s value: %s\n", option.Name, err)
		}
	}
}
func ProcessPushOptions(options *[]Option, pushOptions *git.PushOptions) {
	for _, option := range *options {
		switch option.Name {
		case "RemoteName":
			pushOptions.RemoteName = option.Value
		default:
		}
	}
}
func ProcessPullOptions(options *[]Option, pullOptions *git.PullOptions) {
	var err error
	for _, option := range *options {
		switch option.Name {
		case "SingleBranch":
			pullOptions.SingleBranch, err = strconv.ParseBool(option.Value)
		case "RemoteName":
			pullOptions.RemoteName = option.Value
		case "Depth":
			i, err := strconv.ParseInt(option.Value, 10, 32)
			if err != nil {
				log.Fatalf("Cannot get Depth value: %s\n", err)
			}
			pullOptions.Depth = int(i)
		default:
		}
		if err != nil {
			log.Fatalf("Cannot get %s value: %s\n", option.Name, err)
		}
	}
}
func ProcessCheckoutOptions(options *[]Option, checkoutOptions *git.CheckoutOptions) {
	for _, option := range *options {
		switch option.Name {
		case "Branch", "ReferenceName":
			checkoutOptions.Branch = plumbing.NewBranchReferenceName(option.Value)
		}
	}
}
