package store

import (
	"database/sql"
	"log"

	//Import model ...
	"c/GoExam/imagesUrlColor/model"

	_ "github.com/lib/pq" //...
)

//Store  ...
type Store struct {
	config             *Config
	db                 *sql.DB
	urlImageRepositiry *URLImageRepository
}

//New ...
func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

//Open ...
func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

//Close ...
func (s *Store) Close() {
	s.db.Close()
}

//InsertURL ...
func (r *URLImageRepository) InsertURL(urlColorImages []model.URLImage) {
	for _, u := range urlColorImages {
		r.store.db.Exec("insert into img_Url_Color values($1, $2)", u.URLImg, u.Color)
	}
	log.Print("Insert OK")
}


