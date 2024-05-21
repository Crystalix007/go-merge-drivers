# go-merge-drivers

git merge drivers [1] for go.mod and go.sum files.

Allows seamless editing from multiple members of the team, bumping to the latest
version supported automatically on merge.

## To Use go.mod driver

Define a merge driver like this in `.git/config`:

```gitconfig
[merge "go-mod"]
	name = Go mod merge driver
	driver = go run ./cmd/go-mod-merge/ -O %O -A %A -B %B -P %P --output %A
```

Then define gitattributes that use this driver:

```gitattributes
go.mod merge=go-mod
```
