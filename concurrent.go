package golib

import (
	"sync"
)

type Limiter interface {
	Start()
	End()
	Wait()
}

type ConcurrentLimitWaitGroup struct {
	ch chan struct{}
	wg *sync.WaitGroup
}

func NewConcurrentLimitWaitGroup(concurrent int) Limiter {
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
	<-l.ch
	l.wg.Done()
}

func (l *ConcurrentLimitWaitGroup) Wait() {
	l.wg.Wait()
}

type ConcurrentLimit struct {
	ch chan struct{}
}

func NewConcurrentLimit(concurrent int) Limiter {
	return &ConcurrentLimit{
		ch: make(chan struct{}, concurrent),
	}

}

func (l *ConcurrentLimit) Start() {

	l.ch <- struct{}{}
}

func (l *ConcurrentLimit) End() {
	<-l.ch
}

func (l *ConcurrentLimit) Wait() {
	return
}

type NopLimit struct {
}

func NewNopLimiter() Limiter {
	return &NopLimit{}
}

func (l *NopLimit) Start() {
	return
}
func (l *NopLimit) End() {
	return
}
func (l *NopLimit) Wait() {
	return
}
