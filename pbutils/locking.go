package pbutils

import (
	"os"
	"sync"
	"syscall"
	"fmt"
	"net"
	"strconv"
)

// file lock
type FileLock struct {
	path string
	f    *os.File
	sync.Mutex
}

func NewFLock(path string) (*FileLock, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &FileLock{
		path: path,
		f:    f,
	}, nil
}

func (l *FileLock) TryLock() error {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	if l.f == nil {
		f, err := os.OpenFile(l.path, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		l.f = f
	}

	err := syscall.Flock(int(l.f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return fmt.Errorf("cannot flock path %s - %s", l.path, err)
	}
	return nil
}

func (l *FileLock) WriteContent(s string) error {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	_, err := l.f.WriteString(s)
	return err
}

func (l *FileLock) Unlock() error {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	if l.f == nil {
		return nil
	}
	err := syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
	l.f.Close()
	l.f = nil
	return err
}

// locked by port
type PortLock struct {
	host string
	port int
	ln   net.Listener
}

func NewPortLock(port int) *PortLock {
	return &PortLock{host: "127.0.0.1", port: port}
}

func (p *PortLock) TryLock() (bool, error) {
	if l, err := net.Listen("tcp", net.JoinHostPort(p.host, strconv.Itoa(p.port))); err == nil {
		p.ln = l
		return true, nil
	}
	return false, nil
}

func (p *PortLock) UnLock() error {
	if p.ln == nil {
		return nil
	}
	err := p.ln.Close()
	p.ln = nil
	return err
}

// windows
func (l *FileLock) TryLockWin() error {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	if _, err := os.Stat(l.path); err == nil {
		// If the files exists, we first try to remove it
		if err = os.Remove(l.path); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	file, err := os.OpenFile(l.path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	l.f = file

	return nil
}

func (l *FileLock) UnlockWin() error {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	if l.f == nil {
		return nil
	}
	l.f.Close()
	l.f = nil
	return nil
}
