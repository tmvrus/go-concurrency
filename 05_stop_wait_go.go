package go_concurrency

type storage interface{}

type StorageDB struct{}

func (db *StorageDB) Stop() {}
func NewDB(dsn string) StorageDB {
	return StorageDB{}
}

type API struct{ db storage }

func (a API) Stop() {}
func NewAPI(db storage) API {
	return API{db}
}

type Consumer struct{ db storage }

func (c Consumer) Stop() {}
func NewConsumer(db storage) Consumer {
	return Consumer{db}
}

func Start() {
	db := NewDB("")
	api := NewAPI(db)
	consumer := NewConsumer(db)

	api.Stop()      // API прекращает принимать запросы
	consumer.Stop() // Consumer прекращает принимать сообщения

	db.Stop() // безопасная остановка базы данных
}
