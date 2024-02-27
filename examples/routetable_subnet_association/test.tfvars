routes = {
  "route1" = {
    address_prefix = "0.0.0.0/0",
    next_hop_type  = "Internet"
  },
  "route2" = {
    address_prefix         = "10.0.0.0/8",
    next_hop_type          = "VirtualAppliance",
    next_hop_in_ip_address = "10.10.10.0"
  }
}
logical_product_service = "sbntrtassotn"
network_map = {
  "spoke1" = {
    address_space                                         = ["192.0.0.0/16"]
    subnet_names                                          = ["somesbt"]
    subnet_prefixes                                       = ["192.0.0.0/24"]
    bgp_community                                         = null
    ddos_protection_plan                                  = null
    dns_servers                                           = []
    nsg_ids                                               = {}
    route_tables_ids                                      = {}
    subnet_delegation                                     = {}
    subnet_enforce_private_link_endpoint_network_policies = {}
    subnet_enforce_private_link_service_network_policies  = {}
    subnet_service_endpoints                              = {}
    tags                                                  = {}
    tracing_tags_enabled                                  = false
    tracing_tags_prefix                                   = ""
    use_for_each                                          = true
  }
}
