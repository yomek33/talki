package stores

import (
	"gorm.io/gorm"
)

type Stores struct {
	DB           *gorm.DB
	UserStore    UserStore
	ArticleStore ArticleStore
	PhraseStore  PhraseStore
}

func NewStores(db *gorm.DB) *Stores {
	return &Stores{
		DB:           db,
		UserStore:    &userStore{BaseStore{DB: db}},
		ArticleStore: &articleStore{BaseStore{DB: db}},
		PhraseStore:  &phraseStore{BaseStore{DB: db}},
	}
}

type BaseStore struct {
	DB *gorm.DB
}

func (bs *BaseStore) PerformDBTransaction(fn func(*gorm.DB) error) error {
	tx := bs.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
