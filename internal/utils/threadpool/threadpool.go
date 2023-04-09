package threadpool

import (
	"sync"
)

type ThreadPool struct {
	MaxThreads int
	tasks      chan func()
	wg         sync.WaitGroup
}

func NewThreadPool(maxThreads int) *ThreadPool {
	return &ThreadPool{
		MaxThreads: maxThreads,
		tasks:      make(chan func()),
	}
}

func (p *ThreadPool) AddTask(task func()) {
	p.wg.Add(1)
	p.tasks <- func() {
		task()
		p.wg.Done()
	}
}

func (p *ThreadPool) Start() {
	for i := 0; i < p.MaxThreads; i++ {
		go func() {
			for task := range p.tasks {
				task()
			}
		}()
	}
}

func (p *ThreadPool) Wait() {
	close(p.tasks)
	p.wg.Wait()
}
