package stores

type ChatStore interface {
	CreateChat() error
	GetChatByID() error
}

type chatStore struct {
	BaseStore
}

func (s *chatStore) CreateChat() error {
	return nil
}

func (s *chatStore) GetChatByID() error {
	return nil
}
