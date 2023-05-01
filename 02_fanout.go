package go_concurrency

import (
	"context"
	"sync"
	"time"
)

type Email struct {
	From, To, Title, Body string
}

type db interface {
	Get(ctx context.Context) []Email
}

func StartFanOut(ctx context.Context, source db) {
	pipe := make(chan Email, bufferSize)
	wg := &sync.WaitGroup{}

	for i := 0; i < count; i++ {
		wg.Add(1)
		go emailSender(wg, pipe)
	}

	for {
		tasks := source.Get(ctx) // выбираем большую пачку данных
		if len(tasks) == 0 {
			time.Sleep(time.Second)
			continue
		}

		for i := range tasks { // перкладываем ее в канал
			select {
			case pipe <- tasks[i]:
				continue
			case <-ctx.Done():
				close(pipe)
				break
			}
		}
		break
	}

	wg.Wait() // не выходым пока отправители не вычитают всю пачку данных

}

func emailSender(wg *sync.WaitGroup, source <-chan Email) {
	defer wg.Done()
	for email := range source {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		sendEmail(ctx, email)
		cancel()
	}
}

func sendEmail(ctx context.Context, e Email) {
	// долгий SMTP-вызов
}
