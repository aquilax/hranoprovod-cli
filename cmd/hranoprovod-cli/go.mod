module github.com/aquilax/hranoprovod-cli/cmd/hranoprovod-cli/v3

go 1.17

require (
	github.com/aquilax/hranoprovod-cli/lib/filter/v3 v3.0.0
	github.com/aquilax/hranoprovod-cli/lib/parser/v3 v3.0.0
	github.com/aquilax/hranoprovod-cli/lib/resolver/v3 v3.0.0
	github.com/aquilax/hranoprovod-cli/v3 v3.0.0
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

replace github.com/aquilax/hranoprovod-cli/v3 => ../../
replace github.com/aquilax/hranoprovod-cli/lib/parser/v3 => ../../lib/parser
replace github.com/aquilax/hranoprovod-cli/lib/filter/v3 => ../../lib/filter
replace github.com/aquilax/hranoprovod-cli/lib/resolver/v3 => ../../lib/resolver