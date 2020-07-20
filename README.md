# SQLBundle

`SQLBundle` is a database migration tool. Manage your database schema by creating incremental SQL changes.

# Install

Currently, `sqlbundle` just support `oracle` and `postgresql`, so for a fully version of the binary, use:

```shell
$ make release
```

Please noted that fully version just support Linux. For other version such as `postgresql-only` or `none-db` , binary supports Linux, macOS and Windows



For a lite version of the binary without DB connection dependent commands, use the exclusive build tags:

```shell
$ make dev-nodb #for only current os

#or

$ make release-nodb #for all OS
```



# Usage

```shell
Usage: sqlbundle COMMAND [OPTIONS]

COMMAND:
  init          Init new sql project
  create        Create a new sql file
  install       Download and install dependencies into deps directory
  pack          Packing
  clean         Remove build directory
  publish       Deploy package to repository
  list          List all migrations file
  version       Print version of sql-bundle
  upgrade       Upgrade database to latest version
  downgrade     Downgrade database to previous version or any specific version
  help          Print usage info

Examples:

Options:
  -artifact string
    	artifact name
  -db-connection string
    	connection string of database
  -db-driver string
    	database driver
  -file string
    	name of sql file
  -force
    	force to delete before install
  -group string
    	group name
  -pass string
    	repository credential
  -repo string
    	address of repository
  -user string
    	repository credential
  -verbose
    	print more log
  -version string
    	version of database
  -workdir string
    	working directory
```

## License

Licensed under [GNU Affero General Public License Version 3] [./LICENSE]

