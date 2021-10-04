package generate

//go:generate kratos-cli proto client api ts-umi openapi

//  kratos-cli proto add api/cms/v1/category.proto
//  kratos-cli proto server api/cms/v1/category.proto internal/services/cms/

//  kratos-cli proto add api/sys/v1/sys.proto
//  kratos-cli proto client api/sys/v1/sys.proto
//  kratos-cli proto server api/sys/v1/sys.proto internal/apps/sys/
