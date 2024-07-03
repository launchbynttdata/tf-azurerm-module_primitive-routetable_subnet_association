// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

module "resource_group" {
  source  = "terraform.registry.launch.nttdata.com/module_primitive/resource_group/azurerm"
  version = "~> 1.0"

  name     = local.resource_group_name
  location = var.region
  tags = {
    resource_name = local.resource_group_name
  }
}

# This module generates the resource-name of resources based on resource_type, naming_prefix, env etc.
module "resource_names" {
  source  = "terraform.registry.launch.nttdata.com/module_library/resource_name/launch"
  version = "~> 1.0"

  for_each = var.resource_names_map

  region                  = join("", split("-", var.region))
  class_env               = var.class_env
  cloud_resource_type     = each.value.name
  instance_env            = var.instance_env
  instance_resource       = var.instance_resource
  maximum_length          = each.value.max_length
  logical_product_family  = var.logical_product_family
  logical_product_service = var.logical_product_service
}

module "route_table" {
  source  = "terraform.registry.launch.nttdata.com/module_primitive/route_table/azurerm"
  version = "~> 1.0"

  name                          = local.route_table_name
  location                      = var.region
  resource_group_name           = module.resource_group.name
  disable_bgp_route_propagation = var.disable_bgp_route_propagation
  tags                          = var.tags
  depends_on                    = [module.resource_group]
}

module "network" {
  source  = "terraform.registry.launch.nttdata.com/module_collection/virtual_network/azurerm"
  version = "~> 1.0"

  network_map = local.modified_network_map

  depends_on = [module.resource_group]
}

module "rttbl_subnet_association" {
  source = "../.."

  for_each = module.network.vnet_subnets

  route_table_id = module.route_table.id
  subnet_id      = each.value[0]
}
