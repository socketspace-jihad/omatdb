package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/socketspace-jihad/omatdb/engine"
	omathttp "github.com/socketspace-jihad/omatdb/handler/http"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	db := engine.NewKVStore()
	db.Load()
	go func() {
		httpHandler := omathttp.NewKVHandler(db)
		httpHandler.Run()
	}()

	<-sig
	//doing flush / persisting the data.
	db.Flush()
}
