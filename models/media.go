package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Media struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Type      string  `xorm:"not null varchar(128)"`
	Hash      string  `xorm:"not null varchar(32)"`
	Size      int64   `xorm:"not null"`
	File      string  `xorm:"varchar(256)"`
	Thumbnail *string `xorm:"varchar(256)"`
	Likes     int64   `xorm:"not null"`
	Dislikes  int64   `xorm:"not null"`

	Tags []Tag `gorm:"many2many:m2m_media_tag;"`
}

func (md *Media) New(db *gorm.DB, file, hash, ftype string, fsize int64, thumbnail *string) (*Media, error) {
	md = &Media{
		File:      file,
		Hash:      hash,
		Type:      ftype,
		Size:      fsize,
		Thumbnail: thumbnail,
	}
	return md, db.Create(md).Error
}

func (md *Media) Get(db *gorm.DB, limit int) ([]Media, error) {
	var media []Media
	return media, db.Model(md).Limit(limit).Order("id DESC").Find(&media).Error
}

func (md *Media) GetTo(db *gorm.DB, last_id int, limit int) ([]Media, error) {
	var media []Media
	return media, db.Model(md).Where("id < ?", last_id).Limit(limit).Order("id DESC").Find(&media).Error
}

func (md *Media) GetByHash(db *gorm.DB, hash string) (*Media, bool, error) {
	err := db.Where(&Media{Hash: hash}).First(md).Error
	if err == gorm.ErrRecordNotFound {
		return md, false, nil
	} else {
		return md, true, err
	}
}

func (md *Media) TagsAdd(db *gorm.DB, tag *Tag) error {
	return db.Model(md).Association("Tags").Append(tag).Error
}

func (md *Media) TagsGet(db *gorm.DB) (error, []Tag) {
	var tags []Tag
	return db.Model(md).Association("Tags").Find(&tags).Error, tags
}

func (md *Media) TagsRemove(db *gorm.DB, tag *Tag) error {
	return db.Model(md).Association("Tags").Delete(tag).Error
}
