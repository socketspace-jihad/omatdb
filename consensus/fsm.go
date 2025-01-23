package consensus

import (
	"encoding/json"
	"io"

	"github.com/hashicorp/raft"
)

type fsm Raft

func (f *fsm) Apply(l *raft.Log) interface{} {
	var c Command
	if err := json.Unmarshal(l.Data, &c); err != nil {
		return nil
	}
	switch c.Operation {
	case "post":
		f.KVStorer.Store(c.Key, c.Value)
	case "update":
		f.KVStorer.Update(c.Key, c.Value)
	case "delete":
		f.KVStorer.Delete(c.Key)
	}
	return nil
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (f *fsm) Restore(rc io.ReadCloser) error {
	return nil
}
