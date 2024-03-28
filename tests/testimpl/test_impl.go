package common

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v5"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
)

func TestRouteTableSubnetAssociation(t *testing.T, ctx types.TestContext) {

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		t.Fatal("ARM_SUBSCRIPTION_ID is not set in the environment variables ")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		t.Fatalf("Unable to get credentials: %e\n", err)
	}

	clientFactory, err := armnetwork.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		t.Fatalf("Unable to get clientFactory: %e\n", err)
	}

	subnetsClient := clientFactory.NewSubnetsClient()
	routeTableClient := clientFactory.NewRouteTablesClient()

	resourceGroupName := terraform.Output(t, ctx.TerratestTerraformOptions(), "resource_group_name")
	routeTableName := terraform.Output(t, ctx.TerratestTerraformOptions(), "name")
	vnetNames := terraform.OutputMap(t, ctx.TerratestTerraformOptions(), "vnet_names")
	subnetNames := terraform.OutputMap(t, ctx.TerratestTerraformOptions(), "vnet_subnets")

	t.Run("IsRouteTableSubnetAssociated", func(t *testing.T) {

		routeTable, err := routeTableClient.Get(context.Background(), resourceGroupName, routeTableName, nil)
		if err != nil {
			t.Fatalf("Error getting Route Table: %v", err)
		}
		if routeTable.Name == nil {
			t.Fatalf("Route Table does not exist")
		}

		for _, vnetName := range vnetNames {
			for _, subnetName := range subnetNames {
				inputSubnetName := strings.Trim(getSubstring(subnetName), "[]")

				subnet, err := subnetsClient.Get(context.Background(), resourceGroupName, vnetName, inputSubnetName, nil)
				if err != nil {
					t.Fatalf("Error getting subnet: %v", err)
				}
				if subnet.Name == nil {
					t.Fatalf("Subnet does not exist")
				}
				subnetRouteTable := subnet.Properties.RouteTable
				assert.NotEmpty(t, subnetRouteTable, "Subnet does not have a route table associated.")
			}
		}
	})
}

func getSubstring(input string) string {
	parts := strings.Split(input, "/")
	return parts[len(parts)-1]
}
