package main

import (
  "log"
  "github.com/MindTwister/goon"
  "flag"
  "os/exec"
  "os"
)

func main() {
  var interval int
  var dir string
  flag.IntVar(&interval,"interval",300,"Filesystem check interval in milliseconds (default: 300)")
  flag.StringVar(&dir,"dir",".","Directory to watch (default: .)")
  flag.Parse()
  args := flag.Args()
  notifier := goon.Watch(dir, interval)
  cmdName := args[0]
  for {
    <- notifier
		cmd := exec.Command(cmdName,args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Print("Running:", args)
		err := cmd.Run()
		if err != nil {
			log.Print(err)
		}
		cmd.Wait()
  }
}
