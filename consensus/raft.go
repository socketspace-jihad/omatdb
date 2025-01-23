package consensus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/socketspace-jihad/omatdb/engine"
)

type JoinClusterData struct {
	Addr string `json:"addr"`
	ID   string `json:"id"`
	*Raft
}

func (j *JoinClusterData) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if j.Addr == "" || j.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := j.Raft.Join(j.ID, j.Addr); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func JoinCluster(joinAddr, raftAddr, nodeID string) error {
	b, err := json.Marshal(JoinClusterData{
		Addr: raftAddr,
		ID:   nodeID,
	})
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/join", joinAddr), "application/json", bytes.NewReader(b))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer resp.Body.Close()
	return nil
}

type Raft struct {
	Bind     string
	DataDir  string
	Rft      *raft.Raft
	KVStorer engine.KVStorer
}

func (r *Raft) Open(bootstrapped bool, nodeID string) error {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	addr, err := net.ResolveTCPAddr("tcp", r.Bind)
	if err != nil {
		return err
	}

	transport, err := raft.NewTCPTransport(r.Bind, addr, 3, 20*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	snapshot, err := raft.NewFileSnapshotStore(r.DataDir, 2, os.Stderr)
	if err != nil {
		return err
	}

	logStore := raft.NewInmemStore()
	stableStore := raft.NewInmemStore()

	stateMachine := fsm(Raft{
		KVStorer: r.KVStorer,
	})

	rf, err := raft.NewRaft(
		config,
		&stateMachine,
		logStore, stableStore, snapshot, transport,
	)
	if err != nil {
		return err
	}
	r.Rft = rf

	if bootstrapped {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		if err := rf.BootstrapCluster(configuration); err != nil {
			return err.Error()
		}
	}
	return nil
}

func (r *Raft) Join(nodeID string, sourceAddr string) error {
	configFuture := r.Rft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return err
	}
	for _, server := range configFuture.Configuration().Servers {
		log.Println("Seeing server ID:", server.ID)
		if server.ID == raft.ServerID(nodeID) && server.Address == raft.ServerAddress(sourceAddr) {
			return nil
		}
		if server.ID == raft.ServerID(nodeID) || server.Address == raft.ServerAddress(sourceAddr) {
			future := r.Rft.RemoveServer(server.ID, 0, 0)
			if err := future.Error(); err != nil {
				return err
			}
		}
	}
	f := r.Rft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(sourceAddr), 0, 30*time.Second)
	return f.Error()
}
