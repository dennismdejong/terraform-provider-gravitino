# Complete Gravitino Example

This example demonstrates a full Gravitino deployment with all resource types working together:

- **Metalake** - Root namespace
- **Catalogs** - Hive, Iceberg, Fileset, Kafka, ML
- **Schemas** - Organized by domain
- **Tables** - With columns, partitions, sort order
- **Views** - Analytics views over tables
- **Topics** - Event streaming
- **Filesets** - Data lake storage
- **Functions** - UDFs
- **Models** - ML model registry with versioning
- **Tags** - Data classification
- **Policies** - Access control
- **Users, Groups, Roles** - Identity and authorization
- **Owners** - Resource ownership
- **Jobs** - Scheduled tasks

To use, copy the `provider.tf` from `examples/provider/` alongside this file, or configure your own provider block.
