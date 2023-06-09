module github.com/shamaton/aws-dax-go-v2-test

go 1.20

require (
	github.com/aws/aws-dax-go v1.2.12
	github.com/aws/aws-sdk-go-v2 v1.18.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.10.25
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.19.7
)

require (
	github.com/antlr/antlr4 v0.0.0-20181218183524-be58ebffde8e // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.24 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.33 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.27 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.14.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.27 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.27 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.12.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.14.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.19.0 // indirect
	github.com/aws/smithy-go v1.13.5 // indirect
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

require (
	github.com/aws/aws-sdk-go-v2/config v1.18.25
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.4.52
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
)

replace github.com/aws/aws-dax-go v1.2.12 => github.com/shamaton/aws-dax-go v1.2.12-sdk.v2
