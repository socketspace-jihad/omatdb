![alt text](https://github.com/socketspace-jihad/omatdb/omatdb-logo.png)
# OmatDB - Distributed Key Value Store

## Why i build this ?
Well, i want to learn and implement these things.
- concurrently accessed Map data structure
- how to create in memory first database
- how to persist the data to the disk
- how raft implementation to have a synchronize distributed system
- finite state machine
- how rpc works

# How can i run the database ?

Build from source code
- Pull the source code
	- `git clone https://github.com/socketspace-jihad/omatdb`
- Download dependencies package
	- `go get`
- Compile
	- `go build -o omatdb`

Bootstrap the cluster ( first node )
- `omatdb --bootstrap=true`
---
Available flags
- `bootstrap ( bool, default=false )` determine if the node is for bootstraping new cluster or not
- `dataDir (string, default="./data")` determine where the data would be persisted
- `raftAddr (string, default="127.0.0.1:3300")` determine the listener address for raft consensus protocol
- `httpAddr (string, default="127.0.0.1:8080")` determine the listener address for http server
- `joinAddr (string, default="")` if `bootstrap = false` you have to specify where is the address to join the cluster
- `id (string, default="node1")` determine node id for the cluster ecosystem

---
Spawn new node and joining the cluster
- `omatdb --bootstrap=false --dataDir="./data2" --raftAddr="<should be unique for each node>" --id="<should be unique for each node>" --httpAddr="<should be unique for each node>" --joinAddr="127.0.0.1:3300"`
