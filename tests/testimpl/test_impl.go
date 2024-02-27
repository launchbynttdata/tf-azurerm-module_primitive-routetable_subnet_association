package common

import (
	"context"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/nexient-llc/lcaf-component-terratest-common/lib/azure/configure"
	"github.com/nexient-llc/lcaf-component-terratest-common/lib/azure/login"
	"github.com/nexient-llc/lcaf-component-terratest-common/lib/azure/network"
	"github.com/nexient-llc/lcaf-component-terratest-common/types"
	"github.com/stretchr/testify/assert"
)

const terraformDir string = "../../examples/routetable_subnet_association"
const varFile string = "test.tfvars"

func TestRouteTableSubnetAssociation(t *testing.T, ctx types.TestContext) {

	envVarMap := login.GetEnvironmentVariables()
	clientID := envVarMap["clientID"]
	clientSecret := envVarMap["clientSecret"]
	tenantID := envVarMap["tenantID"]
	subscriptionID := envVarMap["subscriptionID"]

	spt, err := login.GetServicePrincipalToken(clientID, clientSecret, tenantID)
	if err != nil {
		t.Fatalf("Error getting Service Principal Token: %v", err)
	}

	subnetsClient := network.GetSubnetsClient(spt, subscriptionID)
	routeTableClient := network.GetRouteTablesClient(spt, subscriptionID)
	terraformOptions := configure.ConfigureTerraform(terraformDir, []string{terraformDir + "/" + varFile})
	t.Run("IsRouteTableSubnetAssociated", func(t *testing.T) {
		resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
		routeTableName := terraform.Output(t, ctx.TerratestTerraformOptions(), "name")
		vnetNames := terraform.OutputMap(t, ctx.TerratestTerraformOptions(), "vnet_names")
		subnetNames := terraform.OutputMap(t, ctx.TerratestTerraformOptions(), "vnet_subnets")

		routeTable, err := routeTableClient.Get(context.Background(), resourceGroupName, routeTableName, "")
		if err != nil {
			t.Fatalf("Error getting Route Table: %v", err)
		}
		if routeTable.Name == nil {
			t.Fatalf("Route Table does not exist")
		}

		for _, vnetName := range vnetNames {
			for _, subnetName := range subnetNames {
				inputSubnetName := strings.Trim(getSubstring(subnetName), "[]")

				subnet, err := subnetsClient.Get(context.Background(), resourceGroupName, vnetName, inputSubnetName, "")
				if err != nil {
					t.Fatalf("Error getting subnet: %v", err)
				}
				if subnet.Name == nil {
					t.Fatalf("Subnet does not exist")
				}
				subnetRouteTable := subnet.RouteTable
				assert.NotEmpty(t, subnetRouteTable, "Subnet does not have a route table associated.")
			}
		}
	})
}

func getSubstring(input string) string {
	parts := strings.Split(input, "/")
	return parts[len(parts)-1]
}
