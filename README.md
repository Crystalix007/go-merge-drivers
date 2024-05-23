# go-merge-drivers

git merge drivers [[1]] for go.mod and go.sum files.

Allows seamless editing from multiple members of the team, bumping to the latest
version supported automatically on merge.

## To Use go.mod driver

Define a merge driver like this in `.git/config`:

```gitconfig
[merge "go"]
	name = Go merge driver
	driver = go run github.com/Crystalix007/go-merge-drivers/cmd/go-merge/ -O %O -A %A -B %B -P %P --output %A
```

Then define gitattributes that use this driver:

```gitattributes
go.mod merge=go
go.sum merge=go
```

[1]: https://git-scm.com/docs/gitattributes#_defining_a_custom_merge_driver
