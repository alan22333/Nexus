package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name  string  `gorm:"unique;not null;size:50"`
	Posts []*Post `gorm:"many2many:post_tags;"` // 反向关联，方便查询，但通常不在JSON中返回
}
