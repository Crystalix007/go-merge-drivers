module github.com/crystalix007/go-merge-drivers

go 1.24.0

replace github.com/spf13/cobra => gitlab.com/spf13/cobra v1.7.0

tool (
	github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
	github.com/hexdigest/gowrap
)

require github.com/spf13/cobra v1.8.0

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/oapi-codegen/oapi-codegen/v2 v2.4.1 // indirect
	github.com/hexdigest/gowrap v1.4.2 // indirect
	golang.org/x/mod v0.16.0
)
