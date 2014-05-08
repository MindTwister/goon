package goon

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"testing"
	"time"
)

func shouldReceive(notifier chan struct{}) (ok bool) {
	ok = false
	select {
	case <-notifier:
		ok = true
	case <-time.After(time.Second):
	}
	return ok
}

func Test_Watch(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "goon_test")
	genericFileName := path.Join(tmpDir, "genericFile")
	genericFile, err := os.Create(genericFileName)
	os.Mkdir(path.Join(tmpDir, "sub"), os.ModeDir)
	os.Create(path.Join(tmpDir, "sub", "for_show"))
	if err != nil {
		panic(err)
	}
	genericFile.WriteString("test")
	genericFile.Close()
	time.Sleep(time.Millisecond * 100)
	notifier := Watch(tmpDir, 300)
	var waiter sync.WaitGroup
	waiter.Add(4)
	Convey("Given a watched path", t, func() {
		Convey("If a file is created", func() {
			ioutil.TempFile(tmpDir, "")
			Convey("The notifier should have received a message", func() {
				ok := shouldReceive(notifier)
				So(ok, ShouldEqual, true)
				waiter.Done()
			})
		})
		Convey("Or a file changes", func() {
			time.Sleep(time.Second) // Sleep to let the system clock propagate
			os.Chtimes(genericFileName, time.Now(), time.Now())
			Convey("The notifier should have received a message", func() {
				ok := shouldReceive(notifier)
				So(ok, ShouldEqual, true)
				waiter.Done()
			})
		})
		Convey("If no file has changed", func() {
			Convey("No message should have arrived", func() {
				ok := shouldReceive(notifier)
				So(ok, ShouldEqual, false)
				waiter.Done()
			})
		})
		Convey("If a file is removed", func() {
			os.Remove(genericFileName)
			Convey("The notifier should have received a message", func() {
				ok := shouldReceive(notifier)
				So(ok, ShouldEqual, true)
				waiter.Done()
			})
		})
	})
	wgDone := make(chan bool)
	go (func() {
		waiter.Wait()
		wgDone <- true
	})()
	select {
	case <-wgDone:
		// All good
	case <-time.After(10 * time.Second):
		// Something not behaving
		panic("Test too slow, try increasing run time, but be aware that something is amiss")
	}
	os.RemoveAll(tmpDir)
}

func TestKnowsWhenToStop(t *testing.T) {
	Convey("The watch routine should shut down the channel in case of error", t, func() {
		// Create temporary directory
		dir, _ := ioutil.TempDir("", "goon_test")
		// Delete it to cause a panic
		os.RemoveAll(dir)
		notifier := Watch(dir, 300)
		time.Sleep(time.Millisecond*100)
		_, ok := <- notifier
		So(ok,ShouldEqual,false)
	})
}
