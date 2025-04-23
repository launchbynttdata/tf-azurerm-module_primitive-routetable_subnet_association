locals {
  route_table_name    = module.resource_names["route_table"].standard
  resource_group_name = module.resource_names["resource_group"].standard

  override_network_attributes_map = { for vnet_name, vnet in var.network_map : vnet_name => {
    resource_group_name = local.resource_group_name
    vnet_name           = module.resource_names["spoke_vnet"].standard
    location            = var.region
    }
  }

  modified_network_map = {
    for vnet_name, vnet in var.network_map : vnet_name => merge(vnet, local.override_network_attributes_map[vnet_name])
  }

  subnet_map = {
    for item in flatten([
      for network_name, network in module.network.vnet_subnet_name_id_map : [
        for subnet_name, subnet_id in network : {
          key   = subnet_name
          value = subnet_id
        }
      ]
    ]) : item.key => item.value
  }
}
