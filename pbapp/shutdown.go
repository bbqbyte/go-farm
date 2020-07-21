package pbapp

import (
	"os"
	"os/signal"
	"keywea.com/cloud/pblib/pb/events"
	"keywea.com/cloud/pblib/pb/log"
	"sync"
	"syscall"
)

var (
	wg sync.WaitGroup
	cleanupDone = make(chan bool)
)

func Exit(code int) {
	defer func() { recover() }()
	close(cleanupDone)
	os.Exit(code)
}

func Wait() {
	wg.Add(1)

	go func() {
		defer wg.Done()

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan,
			os.Interrupt,
			syscall.SIGHUP,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGKILL,
			syscall.SIGQUIT)

		for {
			var sig os.Signal
			select {
			case sig = <-signalChan:
				switch sig {
				case syscall.SIGHUP:
					plog.Info("Caught SIGHUP. Ignoring")
					continue
				case os.Interrupt:
					plog.Info("Caught SIGINT. Exiting")
				case syscall.SIGTERM:
					plog.Info("Caught SIGTERM. Exiting")
				case syscall.SIGQUIT:
					plog.Info("Caught SIGQUIT. Exiting")
				default:
					plog.Info("Caught signal %v. Exiting", log.Object("signal", sig))
				}
			case <-cleanupDone:
				//case <-time.After(5 * time.Second):
			}
			events.Shutdown()
			return
		}
	}()

	wg.Wait()
}
