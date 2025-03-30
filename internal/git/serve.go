package git

type Storage interface{}

type Server struct {
	config  Config
	storage Storage
}
