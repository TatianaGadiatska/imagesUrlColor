package store

import (
	"c/GoExam/imagesUrlColor/model"
)

//URLImageRepository ...
type URLImageRepository struct {
	store *Store
}

//Create ...
func (r *URLImageRepository) Create(u *model.URLImage) (*model.URLImage, error) {
	return nil, nil
}
