package service

import (
	"chatroom/dao"
	"chatroom/model"
	"chatroom/pkg/filesystem"
	"chatroom/pkg/timeutil"
	"context"
)

type IArticleAnnexService interface {
	Create(ctx context.Context, data *model.ArticleAnnex) error
	UpdateStatus(ctx context.Context, uid int, id int, status int) error
	ForeverDelete(ctx context.Context, uid int, id int) error
}

type ArticleAnnexService struct {
	*dao.Source
	ArticleAnnex *dao.ArticleAnnex
	FileSystem   filesystem.IFilesystem
}

func (s *ArticleAnnexService) Create(ctx context.Context, data *model.ArticleAnnex) error {
	return s.ArticleAnnex.Create(ctx, data)
}

// UpdateStatus 更新附件状态
func (s *ArticleAnnexService) UpdateStatus(ctx context.Context, uid int, id int, status int) error {

	data := map[string]any{
		"status": status,
	}

	if status == 2 {
		data["deleted_at"] = timeutil.DateTime()
	}

	_, err := s.ArticleAnnex.UpdateByWhere(ctx, data, "id = ? and user_id = ?", id, uid)
	return err
}

// ForeverDelete 永久删除笔记附件
func (s *ArticleAnnexService) ForeverDelete(ctx context.Context, uid int, id int) error {

	annex, err := s.ArticleAnnex.FindByWhere(ctx, "id = ? and user_id = ?", id, uid)
	if err != nil {
		return err
	}

	_ = s.FileSystem.Delete(s.FileSystem.BucketPrivateName(), annex.Path)

	return s.Source.Db().Delete(&model.ArticleAnnex{}, id).Error
}
