package go_concurrency

import (
	"fmt"
	"time"
)

const batchSize = 10

type dataSource interface {
	NextBatch(size int) []string
}

func newGenerator(source dataSource) func(timeout time.Duration) (string, error) {
	pipe := make(chan string, bufferSize-1) // в канале на 1 меньше места чем данных в батче
	go func() {
		for {
			// вычитываем данные
			// пмшем их в канал
			// на последней записи "подвисамем"
			// как только ее вычитали, код разблокируется и снова загрузит данные
			batch := source.NextBatch(batchSize)
			for i := range batch {
				pipe <- batch[i]
			}
		}
	}()

	return func(timeout time.Duration) (string, error) {
		// код ждет либо данных из pipe либо сигнала timer
		t := time.After(timeout)
		select {
		case <-t:
			return "", fmt.Errorf("timeout")
		case s := <-pipe:
			return s, nil
		}
	}
}
