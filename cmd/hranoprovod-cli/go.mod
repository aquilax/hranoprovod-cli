module github.com/aquilax/hranoprovod-cli/v2/cmd/hranoprovod-cli

go 1.17

require (
	github.com/aquilax/hranoprovod-cli/v2/lib/filter v0.0.0
	github.com/aquilax/hranoprovod-cli/v2/lib/parser v0.0.0
	github.com/aquilax/hranoprovod-cli/v2/lib/resolver v0.0.0
	github.com/aquilax/hranoprovod-cli/v2/lib/shared v0.0.0
	github.com/stretchr/testify v1.7.1
	github.com/tj/go-naturaldate v1.3.0
	github.com/urfave/cli/v2 v2.3.0
	gopkg.in/gcfg.v1 v1.2.3
)

require (
	github.com/aquilax/truncate v1.0.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

replace github.com/aquilax/hranoprovod-cli/v2/lib/shared => ../../lib/shared

replace github.com/aquilax/hranoprovod-cli/v2/lib/parser => ../../lib/parser

replace github.com/aquilax/hranoprovod-cli/v2/lib/filter => ../../lib/filter

replace github.com/aquilax/hranoprovod-cli/v2/lib/resolver => ../../lib/resolver
