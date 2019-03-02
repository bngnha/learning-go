package videos

import (
	"sync"
)

// ResourceProcessorFunc callback function
type ResourceProcessorFunc func(resource interface{}) (interface{}, error)

// ResultProcessorFunc callback function
type ResultProcessorFunc func(result Result) error

// Job struct
type Job struct {
	id       int
	resource interface{}
}

// Result struct
type Result struct {
	Job   Job
	Extra interface{}
	Err   error
}

// Pool struct
type Pool struct {
	workerNo  int
	jobs      chan Job
	results   chan Result
	done      chan bool
	completed bool
}

// NewPool is a Pool constructor
func NewPool(workerNo int) *Pool {
	r := &Pool{workerNo: workerNo}
	r.jobs = make(chan Job, workerNo)
	r.results = make(chan Result, workerNo)
	r.done = make(chan bool)

	return r
}

func (p *Pool) start(resource []interface{}, procFunc ResourceProcessorFunc, resFunc ResultProcessorFunc) {
	go p.allocate(resource)
	go p.collect(resFunc)
	go p.workerPool(procFunc)

	<-p.done
}

func (p *Pool) allocate(jobs []interface{}) {
	defer close(p.jobs)
	for i, v := range jobs {
		p.jobs <- Job{id: i, resource: v}
	}
}

func (p *Pool) work(wg *sync.WaitGroup, proc ResourceProcessorFunc) {
	defer wg.Done()
	for job := range p.jobs {
		extra, err := proc(job.resource)
		p.results <- Result{job, extra, err}
	}
}

func (p *Pool) workerPool(proc ResourceProcessorFunc) {
	defer close(p.results)
	var wg sync.WaitGroup
	for i := 0; i < p.workerNo; i++ {
		wg.Add(1)
		go p.work(&wg, proc)
	}
	wg.Wait()
}

func (p *Pool) collect(proc ResultProcessorFunc) {
	for result := range p.results {
		proc(result)
	}
	p.done <- true
	p.completed = true
}

// IsCompleted function
func (p *Pool) IsCompleted() bool {
	return p.completed
}
