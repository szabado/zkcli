# zkcli

[![Build Status](https://travis-ci.org/szabado/zkcli.svg?branch=master)](https://travis-ci.org/szabado/zkcli)
[![codecov](https://codecov.io/gh/szabado/zkcli/branch/master/graph/badge.svg)](https://codecov.io/gh/szabado/zkcli)
[![Go Report Card](https://goreportcard.com/badge/github.com/szabado/zkcli)](https://goreportcard.com/report/github.com/szabado/zkcli)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fszabado%2Fzkcli.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fszabado%2Fzkcli?ref=badge_shield)

Elegant, fast, and dependable CLI for ZooKeeper

**zkcli** is a non-interactive command line client for [ZooKeeper](http://zookeeper.apache.org/). It provides with:

 * Basic CRUD-like operations: `create`, `set`, `delete` (aka `rm`), `exists`, `get`, `ls`
 * Extended operations: `lsr` (ls recursive), `creater` (create recursively), `deleter` (aka `rmr`, delete recursively)
 * Well formatted and controlled output: supporting either `txt` or `json` format

### Download & Install

The source code is freely available; you will need `git` installed as well as `go` and `dep`.

### Usage:
```
$ zkcli --help
A CLI to interact with Zookeeper

Usage:
  zkcli [command]

Available Commands:
  create      Create the specified znode
  creater     Create the specified znode, as well as any required parents
  delete      Delete the specified znode
  deleter     Delete the specified znode, as well as any children
  exists      Check if the specified znode exists
  get         Get the value of the specified znode
  getacl      Get the ACL associated with a znode
  help        Help about any command
  ls          Get the children of the specified znode
  lsr         Print the children of the current znode recursively
  set         Set the value of the specified znode
  setacl      Set the ACL of the specified znode

Flags:
      --auth_pwd string   optional, digest scheme, pwd (Can be configured with the environment variable ZKCLI_AUTH_PWD)
      --auth_usr string   optional, digest scheme, user (Can be configured with the environment variable ZKCLI_AUTH_USER)
      --debug             debug mode (very verbose)
      --force             force operation
      --format string     output format (txt|json) (default "txt")
  -h, --help              help for zkcli
      --n                 omit trailing newline
      --servers string    srv1[:port1][,srv2[:port2]...] (Can be configured with the environment variable ZKCLI_SERVERS)
      --verbose           verbose

Use "zkcli [command] --help" for more information about a command.
```

#### Configuration

There are environment variables that can be used instead of command line flags, to reduce the amount that has to be entered

| Flag        | Environment Variable |
|-------------|----------------------|
| `--servers`   | `ZKCLI_SERVERS`        |
| `--auth_user` | `ZKCLI_AUTH_USER`      |
| `--auth_pwd`  | `ZKCLI_AUTH_PWD`       |

### Examples:

```
$ zkcli --servers srv-1,srv-2,srv-3 create /demo_only some_value

# Default port is 2181. The above is equivalent to:
$ zkcli --servers srv-1:2181,srv-2:2181,srv-3:2181 create /demo_only some_value

$ zkcli --servers srv-1,srv-2,srv-3 --format=txt get /demo_only
some_value

# Same as above, JSON format output:
$ zkcli --servers srv-1,srv-2,srv-3 --format=json get /demo_only
"some_value"

# exists exits with exit code 0 when path exists, 1 when path does not exist 
$ zkcli --servers srv-1,srv-2,srv-3 exists /demo_only
true

$ zkcli --servers srv-1,srv-2,srv-3 set /demo_only another_value

$ zkcli --servers srv-1,srv-2,srv-3 --format=json get /demo_only
"another_value"

$ zkcli --servers srv-1,srv-2,srv-3 delete /demo_only

$ zkcli --servers srv-1,srv-2,srv-3 get /demo_only
2014-09-15 04:07:16 FATAL zk: node does not exist

$ zkcli --servers srv-1,srv-2,srv-3 create /demo_only "path placeholder"
$ zkcli --servers srv-1,srv-2,srv-3 create /demo_only/key1 "value1"
$ zkcli --servers srv-1,srv-2,srv-3 create /demo_only/key2 "value2"
$ zkcli --servers srv-1,srv-2,srv-3 create /demo_only/key3 "value3"

$ zkcli --servers srv-1,srv-2,srv-3 ls /demo_only
key3
key2
key1

# Same as above, JSON format output:
$ zkcli --servers srv-1,srv-2,srv-3 --format=json -c ls /demo_only
["key3","key2","key1"]

$ zkcli --servers srv-1,srv-2,srv-3 delete /demo_only
2014-09-15 08:26:31 FATAL zk: node has children

$ zkcli --servers srv-1,srv-2,srv-3 delete /demo_only/key1
$ zkcli --servers srv-1,srv-2,srv-3 delete /demo_only/key2
$ zkcli --servers srv-1,srv-2,srv-3 delete /demo_only/key3
$ zkcli --servers srv-1,srv-2,srv-3 delete /demo_only

# /demo_only path now does not exist.

# Create recursively a path:
$ zkcli --servers=srv-1,srv-2,srv-3 creater "/demo_only/child/key1" "val1"
$ zkcli --servers=srv-1,srv-2,srv-3 creater "/demo_only/child/key2" "val2"

$ zkcli --servers=srv-1,srv-2,srv-3 get "/demo_only/child/key1"
val1

# This path was auto generated due to recursive create:
$ zkcli --servers=srv-1,srv-2,srv-3 get "/demo_only" 
zkcli auto-generated

# ls recursively a path and all sub children:
$ zkcli --servers=srv-1,srv-2,srv-3 lsr "/demo_only" 
child
child/key1
child/key2

# set value with read and write acl using digest authentication
$ zkcli --servers 192.168.59.103 --auth_usr "someuser" --auth_pwd "pass" --acls 1,2 create /secret4 value4

# get value using digest authentication
$ zkcli --servers 192.168.59.103 --auth_usr "someuser" --auth_pwd "pass" get /secret4

# create a value with custom acls
$ zkcli --servers 192.168.59.103 create /secret5 value5 world:anyone:rw,digest:someuser:hashedpw:crdwa

# view the current acl on a path
$ zkcli --servers srv-1,srv-2,srv-3 create /demo_acl "some value"
$ zkcli --servers srv-1,srv-2,srv-3 getacl /demo_acl
world:anyone:cdrwa

# set an acl with world and digest authentication
$ zkcli --servers srv-1,srv-2,srv-3 setacl /demo_acl "world:anyone:rw,digest:someuser:hashedpw:crdwa"
$ zkcli --servers srv-1,srv-2,srv-3 getacl /demo_acl
world:anyone:rw
digest:someuser:hashedpw:cdrwa

# set an acl with world and digest authentication creating the node if it doesn't exist
$ zkcli --servers srv-1,srv-2,srv-3 -force setacl /demo_acl_create "world:anyone:rw,digest:someuser:hashedpw:crdwa"
```

This tool was build because the original [zookeepercli](https://github.com/outbrain/zookeepercli)
has an incredibly verbose interface, very little testing, and
struggles when dealing with recursive operations on large zookeeper instances due to its
single-threaded nature.

**zkcli** aims to be nicer to use, faster, and be more rigorously tested.

### License

Release under the [Apache 2.0 license](https://github.com/szabado/zkcli/blob/master/LICENSE)

Authored by [Felix Jancso-Szabo](https://github.com/szabado) and [Shlomi Noach](https://github.com/shlomi-noach), among others.
 
 
 
 

 


[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fszabado%2Fzkcli.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fszabado%2Fzkcli?ref=badge_large)