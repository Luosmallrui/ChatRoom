package dao

import (
	"chatroom/model"
	"context"

	"gorm.io/gorm"
)

type FileUpload struct {
	Repo[model.FileUpload]
}

func NewFileUpload(db *gorm.DB) *FileUpload {
	return &FileUpload{Repo: NewRepo[model.FileUpload](db)}
}

func (s *FileUpload) GetSplitList(ctx context.Context, uploadId string) ([]*model.FileUpload, error) {
	return s.Repo.FindAll(ctx, func(db *gorm.DB) {
		db.Where("upload_id = ? and type = 2", uploadId)
	})
}

func (s *FileUpload) GetFile(ctx context.Context, uid int, uploadId string) (*model.FileUpload, error) {
	return s.Repo.FindByWhere(ctx, "user_id = ? and upload_id = ? and type = 1", uid, uploadId)
}
