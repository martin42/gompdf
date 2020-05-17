package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mazzegi/gompdf"
)

func main() {
	source := flag.String("source", "../../samples/doc2.xml", "")
	target := flag.String("target", "doc2.pdf", "")
	flag.Parse()

	logf("run watcher (source:%s) (target:%s)", *source, *target)
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	err = watch.Add(*source)
	if err != nil {
		panic(err)
	}
	exec := func() {
		logf("parse and build ...")
		err := gompdf.ParseAndBuildFile(*source, *target)
		if err != nil {
			logf("ERROR: %v", err)
			return
		}
		logf("parse and build ... done")
	}
	notC := make(chan os.Signal)
	signal.Notify(notC, os.Interrupt, os.Kill)
	go func() {
		select {
		case <-notC:
			logf("stop watcher")
			watch.Close()
			return
		}
	}()
	exec()
	lastEvtTime := time.Now()
	for evt := range watch.Events {
		if evt.Op == fsnotify.Write && time.Since(lastEvtTime) > 200*time.Millisecond {
			time.AfterFunc(100*time.Millisecond, exec)
			lastEvtTime = time.Now()
		}
	}
	logf("watcher stopped")
}

func logf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
