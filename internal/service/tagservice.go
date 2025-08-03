package service

import (
	"Nuxus/internal/dao"
	"Nuxus/internal/dto"
	"Nuxus/pkg/erru"
)

func ListTags(sortedBy string) ([]*dto.ListTagsResDTO, error) {

	listTags, err := dao.ListTags(sortedBy)
	if err != nil {
		return nil, erru.ErrInternalServer.Wrap(err)
	}
	return listTags, nil
}
