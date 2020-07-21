package events

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	handlers = []func(){}

	once sync.Once

	wg sync.WaitGroup
	cleanupDone = make(chan bool)
)

func Exit(code int) {
	runHandlers()
	close(cleanupDone)
	os.Exit(code)
}

// add at tail
func RegisterExitHandlerTail(handler func()) {
	handlers = append(handlers, handler)
}

// add at head
func RegisterExitHandlerHead(handler func()) {
	handlers = append([]func(){handler}, handlers...)
}

func runHandler(handler func()) {
	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error: exit handler error:", err)
		}
	}()

	handler()
}

func runHandlers() {
	once.Do(func() {
		for _, handler := range handlers {
			runHandler(handler)
		}
	})
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
					fmt.Println("Caught SIGHUP. Ignoring")
					continue
				case os.Interrupt:
					fmt.Println("Caught SIGINT. Exiting")
				case syscall.SIGTERM:
					fmt.Println("Caught SIGTERM. Exiting")
				case syscall.SIGQUIT:
					fmt.Println("Caught SIGQUIT. Exiting")
				default:
					fmt.Println("Caught signal %v. Exiting", sig)
				}
			case <-cleanupDone:
				//case <-time.After(5 * time.Second):
			}
			return
		}
	}()

	wg.Wait()
}
