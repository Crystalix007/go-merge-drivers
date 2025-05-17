module github.com/crystalix007/go-merge-drivers

go 1.22.2

replace github.com/spf13/cobra => gitlab.com/spf13/cobra v1.8.0

tool (
	github.com/dmarkham/enumer
	github.com/hexdigest/gowrap
)

require github.com/spf13/cobra v1.7.0

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/hexdigest/gowrap v1.4.1 // indirect
	github.com/spf13/pflag v1.0.1
	golang.org/x/mod v0.17.0
)
