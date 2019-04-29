module github.com/seamounts/go-mod-tidy

go 1.12

replace (
	golang.org/x/crypto v0.0.0-20181203042331-505ab145d0a9 => github.com/golang/crypto v0.0.0-20181203042331-505ab145d0a9
	golang.org/x/sys v0.0.0-20181205085412-a5c9d58dba9a => github.com/golang/sys v0.0.0-20181205085412-a5c9d58dba9a
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
)

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.2 // indirect
)
