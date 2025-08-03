package dao

import (
	"Nuxus/internal/dto"
	"Nuxus/internal/models"
)

func ListTags(sortedBy string) ([]*dto.ListTagsResDTO, error) {
	var results []*dto.ListTagsResDTO

	query := DB.Model(&models.Tag{}).
		Select("tags.id, tags.name, count(post_tags.tag_id) as post_count").
		Joins("LEFT JOIN post_tags ON tags.id = post_tags.tag_id").
		Group("tags.id, tags.name")

	switch sortedBy {
	case "post_count":
		query = query.Order("post_count DESC")
	case "name":
		query = query.Order("tags.name ASC")
	default:
		query = query.Order("post_count DESC")
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func FindOrCreateTagByName(name string) (*models.Tag, error) {
	var tag models.Tag
	if err := DB.Where(models.Tag{Name: name}).FirstOrCreate(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func FindOrCreateTagsByNames(names []string) ([]*models.Tag, error) {
	tags := make([]*models.Tag, 0, len(names))
	for _, name := range names {
		tag, err := FindOrCreateTagByName(name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
