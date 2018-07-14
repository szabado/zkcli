# zkcli

[![codecov](https://codecov.io/gh/fJancsoSzabo/zkcli/branch/master/graph/badge.svg)](https://codecov.io/gh/fJancsoSzabo/zkcli)

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
  create      
  creater     
  delete      
  deleter     
  exists      
  get         
  getacl      
  help        Help about any command
  ls          
  lsr         
  rm          
  rmr         
  set         
  setacl      

Flags:
      --auth_pwd string   optional, digest scheme, pwd
      --auth_usr string   optional, digest scheme, user
      --debug             debug mode (very verbose)
      --force             force operation
      --format string     output format (txt|json) (default "txt")
  -h, --help              help for zkcli
      --n                 omit trailing newline
      --servers string    srv1[:port1][,srv2[:port2]...]
      --verbose           verbose
```

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
has an incredibly verbose interface interface, very little testing, ignored any pull requests, and
struggles when dealing with recursive operations on large zookeeper instances due to its
single-threaded nature.

**zkcli** aims to be nicer to use, faster, and be more rigorously tested, as well as be more friendly to
pull and feature requests.

### License

Release under the [Apache 2.0 license](https://github.com/fJancsoSzabo/zkcli/blob/master/LICENSE)

Authored by [Felix Jancso-Szabo](https://github.com/fJancsoSzabo) and [Shlomi Noach](https://github.com/shlomi-noach), among others.
 
 
 
 

 
