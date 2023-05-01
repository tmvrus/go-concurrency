package go_concurrency

import (
	"context"
	"sync"
	"time"
)

type CounterMetric struct {
	ID    string
	Value int
}

const (
	count      = 10
	bufferSize = 1000
)

func StartFanIn(ctx context.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan CounterMetric, bufferSize)
	wgWriter := &sync.WaitGroup{}
	for i := 0; i < count; i++ {
		wgWriter.Add(1)
		go writer(ctx, wgWriter, in)
	}

	out := make(chan CounterMetric)

	wgReader := &sync.WaitGroup{}
	wgReader.Add(1)
	go fainInReader(wgReader, in, out)

	cancel()

	wgWriter.Wait() // ждем остановки всех писателей

	close(in) // читатель получает сигнал что надо завершатся

	wgReader.Wait() // ждем остановки читателя
}

func fainInReader(wg *sync.WaitGroup, in <-chan CounterMetric, out chan<- CounterMetric) {
	defer wg.Done()
	data := make(map[string]int)

	t := time.NewTicker(time.Second)
	for {
		select {
		case m, ok := <-in:
			if !ok {
				flush(data, out)
				return
			}
			data[m.ID] += m.Value
		case <-t.C:
			if len(data) > 0 {
				flush(data, out)
				data = make(map[string]int)
			}
		}
	}
}

func flush(m map[string]int, out chan<- CounterMetric) {
	for id, value := range m {
		out <- CounterMetric{
			ID:    id,
			Value: value,
		}
	}
}

func writer(ctx context.Context, wg *sync.WaitGroup, in chan<- CounterMetric) {
	wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case in <- CounterMetric{}:
			continue
		}
	}
}
