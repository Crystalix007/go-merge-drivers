module github.com/crystalix007/go-merge-drivers

go 1.24.0

require (
	github.com/spf13/cobra v1.8.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/mod v0.17.0
)

require (
	github.com/hexdigest/gowrap v1.4.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/oapi-codegen/oapi-codegen/v2 v2.4.1 // indirect
)

replace github.com/spf13/cobra => gitlab.com/spf13/cobra v1.8.0

tool (
	github.com/dmarkham/enumer
	github.com/hexdigest/gowrap
	github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
)
