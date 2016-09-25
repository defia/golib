package golib

import (
	"sync"
)

type ConcurrentLimitWaitGroup struct {
	ch chan struct{}
	wg *sync.WaitGroup
}

func NewConcurrentLimitWaitGroup(concurrent int) *ConcurrentLimitWaitGroup {
	return &ConcurrentLimitWaitGroup{
		ch: make(chan struct{}, concurrent),
		wg: &sync.WaitGroup{},
	}

}

func (l *ConcurrentLimitWaitGroup) Start() {
	l.wg.Add(1)
	l.ch <- struct{}{}
}

func (l *ConcurrentLimitWaitGroup) End() {
	l.wg.Done()
	<-l.ch
}

func (l *ConcurrentLimitWaitGroup) Wait() {
	l.wg.Wait()
}
