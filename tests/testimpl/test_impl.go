package common

import (
	"context"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v5"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
)

func TestComposableRouteTableSubnetAssociation(t *testing.T, ctx types.TestContext) {

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

	resourceGroupName := terraform.OutputContext(t, context.Background(), ctx.TerratestTerraformOptions(), "resource_group_name")
	routeTableName := terraform.OutputContext(t, context.Background(), ctx.TerratestTerraformOptions(), "name")
	subnetIDs := terraform.OutputMapContext(t, context.Background(), ctx.TerratestTerraformOptions(), "subnet_ids")

	t.Run("IsRouteTableSubnetAssociated", func(t *testing.T) {

		routeTable, err := routeTableClient.Get(context.Background(), resourceGroupName, routeTableName, nil)
		if err != nil {
			t.Fatalf("Error getting Route Table: %v", err)
		}
		if routeTable.Name == nil {
			t.Fatalf("Route Table does not exist")
		}
		if routeTable.ID == nil {
			t.Fatalf("Route Table ID is nil")
		}
		expectedRouteTableID := *routeTable.ID

		for _, subnetID := range subnetIDs {
			parsedSubnetID, err := arm.ParseResourceID(subnetID)
			if err != nil {
				t.Fatalf("Error parsing subnet ID %q: %v", subnetID, err)
			}

			subnet, err := subnetsClient.Get(
				context.Background(),
				parsedSubnetID.ResourceGroupName,
				parsedSubnetID.Parent.Name,
				parsedSubnetID.Name,
				nil,
			)
			if err != nil {
				t.Fatalf("Error getting subnet: %v", err)
			}
			if subnet.Name == nil {
				t.Fatalf("Subnet does not exist")
			}
			assert.NotNil(t, subnet.Properties.RouteTable, "Subnet does not have a route table associated.")
			assert.NotNil(t, subnet.Properties.RouteTable.ID, "Subnet route table ID is nil.")
			assert.Equal(t, expectedRouteTableID, *subnet.Properties.RouteTable.ID, "Subnet is not associated with the expected route table.")
		}
	})
}
