---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pinot_instances Data Source - terraform-provider-pinot"
subcategory: ""
description: |-
  
---

# pinot_instances (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `instances` (Attributes List) The list of instances. (see [below for nested schema](#nestedatt--instances))

<a id="nestedatt--instances"></a>
### Nested Schema for `instances`

Read-Only:

- `admin_port` (Number) The admin port of the instance.
- `enabled` (Boolean) If the instance is enabled.
- `grpc_port` (Number) The GRPC port of the instance.
- `host_name` (String) The hostname of the instance.
- `instance_name` (String) The name of the instance.
- `pools` (List of String) The list of pools.
- `port` (String) The port of the instance.
- `query_mailbox_port` (Number) The query mailbox port of the instance.
- `query_service_port` (Number) The query server port of the instance.
- `system_resource_info` (Attributes) The role of the user. (see [below for nested schema](#nestedatt--instances--system_resource_info))
- `tags` (List of String) The list of tags.

<a id="nestedatt--instances--system_resource_info"></a>
### Nested Schema for `instances.system_resource_info`

Read-Only:

- `max_heap_size_mb` (String) The max heap size in MB.
- `num_cores` (String) The number of cores.
- `total_memory_mb` (String) The total memory in MB.
