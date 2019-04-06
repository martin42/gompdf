package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/martin42/gompdf"
)

func main() {
	source := flag.String("source", "../../samples/doc1.xml", "")
	target := flag.String("target", "doc1.pdf", "")
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
		err := parseAndBuild(*source, *target)
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
	for evt := range watch.Events {
		if evt.Op == fsnotify.Write {
			time.AfterFunc(50*time.Millisecond, exec)
		}
	}
	logf("watcher stopped")
}

func parseAndBuild(source string, target string) error {
	doc, err := gompdf.LoadFromFile(source)
	if err != nil {
		return err
	}
	logf("doc: %v", doc)

	outF, err := os.Create(target)
	if err != nil {
		return err
	}
	defer outF.Close()

	p, err := gompdf.NewProcessor(doc)
	if err != nil {
		return err
	}
	err = p.Process(outF)
	if err != nil {
		return err
	}
	return nil
}

func logf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
