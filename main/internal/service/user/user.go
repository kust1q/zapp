package user

type service struct {
	db       db
	media    mediaService
	producer eventProducer
}

func NewUserService(db db, media mediaService, producer eventProducer) *service {
	return &service{
		db:       db,
		media:    media,
		producer: producer,
	}
}
