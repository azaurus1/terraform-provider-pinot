---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pinot_schemas Data Source - terraform-provider-pinot"
subcategory: ""
description: |-
  
---

# pinot_schemas (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `schemas` (Attributes List) (see [below for nested schema](#nestedatt--schemas))

<a id="nestedatt--schemas"></a>
### Nested Schema for `schemas`

Required:

- `schema_name` (String) The name of the schema.

Optional:

- `date_time_field_specs` (Attributes List) The dimension field specs. (see [below for nested schema](#nestedatt--schemas--date_time_field_specs))
- `dimension_field_specs` (Attributes List) The dimension field specs. (see [below for nested schema](#nestedatt--schemas--dimension_field_specs))
- `enable_column_based_null_handling` (Boolean) Whether to enable column based null handling.
- `metric_field_specs` (Attributes List) The dimension field specs. (see [below for nested schema](#nestedatt--schemas--metric_field_specs))
- `primary_key_columns` (List of String) The primary key columns.

<a id="nestedatt--schemas--date_time_field_specs"></a>
### Nested Schema for `schemas.date_time_field_specs`

Required:

- `data_type` (String) The data type of the dimension.
- `name` (String) The name of the dimension.

Optional:

- `format` (String) The format of the date time.
- `granularity` (String) The granularity of the date time.
- `not_null` (Boolean) Whether the dimension is not null.
- `transform_function` (String) Transform function for specific field.


<a id="nestedatt--schemas--dimension_field_specs"></a>
### Nested Schema for `schemas.dimension_field_specs`

Required:

- `data_type` (String) The data type of the dimension.
- `name` (String) The name of the dimension.

Optional:

- `not_null` (Boolean) Whether the dimension is not null.
- `single_value_field` (Boolean) Whether the dimension is a single value field.
- `transform_function` (String) Transform function for specific field.


<a id="nestedatt--schemas--metric_field_specs"></a>
### Nested Schema for `schemas.metric_field_specs`

Required:

- `data_type` (String) The data type of the dimension.
- `name` (String) The name of the dimension.

Optional:

- `not_null` (Boolean) Whether the dimension is not null.
- `transform_function` (String) Transform function for specific field.
