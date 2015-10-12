# etcdsh
A shell for etcd

This is a shell for interacting with the hierarchy of keys and values withing an etcd cluster by navigating it as if it were a filesystem.

On start up, you start with the present working directory (pwd) set to "/"

Command line flags:

--peers (or -p) (comma delimited list of machine URLS to connect to)
--username (or -u) (username[:password] - if no password, one will be prompted for)
-d (debug)

Commands:

ls, pwd, set, get, env, exit

### pwd

Print present working directory

### ls

List content of directory

ls - list contents of pwd
ls _wildcard_ - list contents of pwd that match the wildcard
ls _path_ - list contents of path 
ls _path_/_wildcard_ - list contents of path that match the wildcard
(wild card paths not currently implemented)

### set

Sets value of key

set _path_ _value_

(value can be quoted)

### get

Prints value of key

get _path_

### env

Set an internal env variable

* e[nv] s[imple] - Print values of keys only
* e[nv] j[son] - Print the result as json
* e[nv] p[retty] - Print the result as formatted json

