package models

const (
	PrivilegeCreateCatalog    = "CREATE_CATALOG"
	PrivilegeUseCatalog       = "USE_CATALOG"
	PrivilegeCreateSchema     = "CREATE_SCHEMA"
	PrivilegeUseSchema        = "USE_SCHEMA"
	PrivilegeCreateTable      = "CREATE_TABLE"
	PrivilegeModifyTable      = "MODIFY_TABLE"
	PrivilegeSelectTable      = "SELECT_TABLE"
	PrivilegeCreateFileset    = "CREATE_FILESET"
	PrivilegeWriteFileset     = "WRITE_FILESET"
	PrivilegeReadFileset      = "READ_FILESET"
	PrivilegeCreateTopic      = "CREATE_TOPIC"
	PrivilegeProduceTopic     = "PRODUCE_TOPIC"
	PrivilegeConsumeTopic     = "CONSUME_TOPIC"
	PrivilegeManageUsers      = "MANAGE_USERS"
	PrivilegeManageGroups     = "MANAGE_GROUPS"
	PrivilegeCreateRole       = "CREATE_ROLE"
	PrivilegeManageGrants     = "MANAGE_GRANTS"
	PrivilegeRegisterModel    = "REGISTER_MODEL"
	PrivilegeLinkModelVersion = "LINK_MODEL_VERSION"
	PrivilegeUseModel         = "USE_MODEL"
	PrivilegeRegisterFunction = "REGISTER_FUNCTION"
	PrivilegeExecuteFunction  = "EXECUTE_FUNCTION"
	PrivilegeModifyFunction   = "MODIFY_FUNCTION"
	PrivilegeCreateTag        = "CREATE_TAG"
	PrivilegeApplyTag         = "APPLY_TAG"
	PrivilegeCreatePolicy     = "CREATE_POLICY"
	PrivilegeApplyPolicy      = "APPLY_POLICY"
	PrivilegeRegisterJob      = "REGISTER_JOB_TEMPLATE"
	PrivilegeUseJob           = "USE_JOB_TEMPLATE"
	PrivilegeRunJob           = "RUN_JOB"
)

var AllPrivileges = []string{
	PrivilegeCreateCatalog,
	PrivilegeUseCatalog,
	PrivilegeCreateSchema,
	PrivilegeUseSchema,
	PrivilegeCreateTable,
	PrivilegeModifyTable,
	PrivilegeSelectTable,
	PrivilegeCreateFileset,
	PrivilegeWriteFileset,
	PrivilegeReadFileset,
	PrivilegeCreateTopic,
	PrivilegeProduceTopic,
	PrivilegeConsumeTopic,
	PrivilegeManageUsers,
	PrivilegeManageGroups,
	PrivilegeCreateRole,
	PrivilegeManageGrants,
	PrivilegeRegisterModel,
	PrivilegeLinkModelVersion,
	PrivilegeUseModel,
	PrivilegeRegisterFunction,
	PrivilegeExecuteFunction,
	PrivilegeModifyFunction,
	PrivilegeCreateTag,
	PrivilegeApplyTag,
	PrivilegeCreatePolicy,
	PrivilegeApplyPolicy,
	PrivilegeRegisterJob,
	PrivilegeUseJob,
	PrivilegeRunJob,
}

const (
	ObjectTypeMetalake    = "METALAKE"
	ObjectTypeCatalog     = "CATALOG"
	ObjectTypeSchema      = "SCHEMA"
	ObjectTypeTable       = "TABLE"
	ObjectTypeColumn      = "COLUMN"
	ObjectTypeFileset     = "FILESET"
	ObjectTypeTopic       = "TOPIC"
	ObjectTypeRole        = "ROLE"
	ObjectTypeModel       = "MODEL"
	ObjectTypeFunction    = "FUNCTION"
	ObjectTypeTag         = "TAG"
	ObjectTypePolicy      = "POLICY"
	ObjectTypeJobTemplate = "JOB_TEMPLATE"
)

var AllObjectTypes = []string{
	ObjectTypeMetalake,
	ObjectTypeCatalog,
	ObjectTypeSchema,
	ObjectTypeTable,
	ObjectTypeFileset,
	ObjectTypeTopic,
	ObjectTypeRole,
	ObjectTypeModel,
	ObjectTypeFunction,
	ObjectTypeTag,
	ObjectTypePolicy,
	ObjectTypeJobTemplate,
}

const (
	PrivilegeConditionAllow = "ALLOW"
	PrivilegeConditionDeny  = "DENY"
)

const (
	OwnerTypeUser  = "USER"
	OwnerTypeGroup = "GROUP"
)

var OwnerObjectTypes = []string{
	ObjectTypeMetalake,
	ObjectTypeCatalog,
	ObjectTypeSchema,
	ObjectTypeTable,
	ObjectTypeFileset,
	ObjectTypeTopic,
	ObjectTypeRole,
}

var StatisticsObjectTypes = []string{
	ObjectTypeMetalake,
	ObjectTypeCatalog,
	ObjectTypeSchema,
	ObjectTypeTable,
	ObjectTypeColumn,
	ObjectTypeFileset,
	ObjectTypeTopic,
	ObjectTypeModel,
	ObjectTypeRole,
}
