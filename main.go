package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/socketspace-jihad/omatdb/consensus"
	"github.com/socketspace-jihad/omatdb/engine"
	omathttp "github.com/socketspace-jihad/omatdb/handler/http"
)

var (
	bootstrapped                                  bool
	dataDir, raftAddr, httpAddr, joinAddr, nodeID string
)

func init() {
	flag.BoolVar(&bootstrapped, "bootstrap", true, "to determine if the startup action is bootstraping the cluster or joining the cluster")
	flag.StringVar(&dataDir, "dataDir", "./data", "to determine where is omatdb persist the data")
	flag.StringVar(&raftAddr, "raftAddr", "127.0.0.1:3300", "to determine rpc handler for raft consensus")
	flag.StringVar(&httpAddr, "httpAddr", "127.0.0.1:8080", "to determine http handler for receiving client")
	flag.StringVar(&joinAddr, "joinAddr", "", "to determine cluster address where the node will be joined")
	flag.StringVar(&nodeID, "id", "node1", "to determine node id inside the cluster of raft consensus ecosystem")
}

func validate() error {
	if !bootstrapped && joinAddr == "" {
		return errors.New("as a non-bootstrap node, you need to specify where this node should join the cluster")
	}
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()

	if err := validate(); err != nil {
		log.Fatalln(err.Error())
	}

	db := engine.NewKVStore(dataDir)
	cnss := consensus.Raft{
		Bind:     raftAddr,
		DataDir:  dataDir,
		KVStorer: db,
	}

	if err := cnss.Open(bootstrapped, nodeID); err != nil {
		log.Fatalln(err.Error())
	}

	if !bootstrapped {
		if err := consensus.JoinCluster(joinAddr, raftAddr, nodeID); err != nil {
			log.Fatalln(err.Error())
		}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	db.Load()
	go func() {
		httpHandler := omathttp.NewKVHandler(db)
		jc := consensus.JoinClusterData{
			Raft: &cnss,
		}
		httpHandler.ServeMux.Handle("/join", &jc)
		httpHandler.Run(httpAddr, &cnss)
	}()

	<-sig
	//doing flush / persisting the data.
	db.Flush()
}
