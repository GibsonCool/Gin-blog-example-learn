package models

import (
	"github.com/jinzhu/gorm"
)

type Article struct {
	Model

	TagId int `json:"tag_id" gorm:"index"` //gorm:index，用于声明这个字段为索引，如果你使用了自动迁移功能则会有所影响，在不使用则无影响
	Tag   Tag `json:"tag"`

	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	State         int    `json:"state"`
	CoverImageUrl string `json:"cover_image_url"`
}

// 在 models.go 中使用自定义 callbacks 就不用每个文件重复写这种回调了
//func (*Article) BeforeCreate(scope *gorm.Scope) error {
//	_ = scope.SetColumn("CreatedOn", time.Now().Unix())
//	return nil
//}
//
//func (*Article) BeforeUpdate(scope *gorm.Scope) error {
//	_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
//	return nil
//}

// ExistArticleByID 通过 id 检测文章是或否已经存在
func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if article.ID > 0 {
		return true, nil
	}
	return false, nil
}

// GetArticleTotal 获取文章的总数
func GetArticleTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Article{}).Where(maps).Count(&count).Error; err != nil {
		return 0, nil
	}

	return count, nil
}

// GetArticleList 根据分页条件获取文章列表信息
func GetArticleList(pageNum int, pageSize int, maps interface{}) ([]*Article, error) {
	var articles []*Article
	err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return articles, nil
}

// GetArticle 通过 id 获取单个文章信息
func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ?", id).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	err = db.Model(&article).Related(&article.Tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &article, nil
}

// EditArticle 根据 id 编辑文章信息
func EditArticle(id int, data interface{}) error {
	if err := db.Model(&Article{}).Where("id = ?", id).Update(data).Error; err != nil {
		return err
	}
	return nil
}

// AddArticle 新增一个文章
func AddArticle(data map[string]interface{}) error {
	err := db.Create(&Article{
		TagId:         data["tag_id"].(int),
		Title:         data["title"].(string),
		Desc:          data["desc"].(string),
		Content:       data["content"].(string),
		CreatedBy:     data["created_by"].(string),
		State:         data["state"].(int),
		CoverImageUrl: data["cover_image_url"].(string),
	}).Error

	if err != nil {
		return err
	}

	return nil
}

// DeleteArticle 通过 id 删除指定文章
func DeleteArticle(id int) error {
	if err := db.Where("id = ?", id).Delete(&Article{}).Error; err != nil {
		return err
	}
	return nil
}

// 硬删除标签
func CleanAllArticle() error {
	// 使用 Unscoped 查找软删除的记录并执行硬删除
	if err := db.Unscoped().Where("deleted_on != ?", 0).Delete(&Article{}).Error; err != nil {
		return err
	}
	return nil
}
