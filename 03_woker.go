package go_concurrency

import "fmt"

type Service struct {
	notify chan int
}

func (s Service) CreateUser(id int, name string) {
	// сохранение пользователя в БД

	s.notify <- id // нотификация о том, что пользователя создан

}

func (s Service) startWorker() {
	for id := range s.notify {
		sendAsyncNotification(id)
	}
}

func NewService() Service {
	s := Service{ // создаем сервис
		notify: make(chan int, bufferSize),
	}

	go s.startWorker() // запускаем воркер

	return s // возвращаем сервиса, вызывающая сторона ничего не знает про воркер
}

func sendAsyncNotification(id int) {
	message := fmt.Sprintf("User with ID %d has been created.", id)
	_ = message // отправка сообщения
}
