package model

import (
	"gorm.io/gorm"
	"time"
)

type Single struct {
	gorm.Model
	Name        string     `gorm:"column:name" json:"Name"`                //type:string       comment:单曲名        version:2023-07-28 23:48
	Intro       string     `gorm:"column:intro" json:"Intro"`              //type:string       comment:介绍          version:2023-07-28 23:48
	Tag         string     `gorm:"column:tag" json:"Tag"`                  //type:string       comment:标签          version:2023-07-28 23:48
	PublishDate *time.Time `gorm:"column:publish_date" json:"PublishDate"` //type:*time.Time   comment:发行日期      version:2023-07-28 23:48
	UserId      uint       `gorm:"column:user_id" json:"UserId"`           //type:*int         comment:发布者id      version:2023-07-28 23:48
	CoverId     uint       `gorm:"column:cover_id" json:"CoverId"`         //type:string       comment:封面物料id    version:2023-07-28 23:48
	DemoId      uint       `gorm:"column:demo_id" json:"DemoId"`           //type:string       comment:预览小样id    version:2023-07-28 23:48
	Tracks      []Track    `gorm:"foreignKey:SingleId" json:"Tracks"`
}

func (s *Single) TableName() string {
	return "single"
}

// Track 分轨
type Track struct {
	gorm.Model
	SingleId     uint    `gorm:"column:single_id" json:"SingleId"`         //type:*int         comment:单曲id      version:2023-07-29 17:05
	Position     string  `gorm:"column:position" json:"Position"`          //type:string       comment:轨道编号    version:2023-07-29 17:05
	Label        string  `gorm:"column:label" json:"Label"`                //type:string       comment:轨道介绍    version:2023-07-29 17:05
	ReleaseQty   uint    `gorm:"column:release_qty" json:"ReleaseQty"`     //type:*int         comment:发行数量    version:2023-07-29 17:05
	ReleasePrice float64 `gorm:"column:release_price" json:"ReleasePrice"` //type:*float64     comment:发行价格    version:2023-07-29 17:05
	InStockQty   uint    `gorm:"column:inStock_qty" json:"InStockQty"`     //type:*int         comment:可售数量    version:2023-07-29 17:05
	FileId       uint    `gorm:"column:file_id" json:"-"`
}

func (s *Track) TableName() string {
	return "track"
}

// File s3文件对象
type File struct {
	gorm.Model
	Type     string `gorm:"column:type" json:"type"`
	Url      string `gorm:"column:url" json:"url"`
	FileName string `gorm:"column:file_name" json:"fileName"`
}

func (s *File) TableName() string {
	return "file"
}
