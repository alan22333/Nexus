package service

import (
	"Nuxus/internal/dao"
	"Nuxus/internal/dto"
	"Nuxus/pkg/erru"
)

type TagService struct {
	tagDAO *dao.TagDAO
}

func NewTagService(tagDAO *dao.TagDAO) *TagService{
	return &TagService{tagDAO: tagDAO}
}

func (t *TagService) ListTags(sortedBy string) ([]*dto.ListTagsResDTO, error) {

	listTags, err := t.tagDAO.ListTags(sortedBy)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}
	return listTags, nil
}
