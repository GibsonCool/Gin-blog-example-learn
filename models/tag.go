package models

import "github.com/jinzhu/gorm"

// 创建 Tag struct 用于 Gorm 使用。并给予附属属性 json。便于在接口返回数据的时候自动转换格式
type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

/*
	这属于gorm的Callbacks，可以将回调方法定义为模型结构的指针，
	在创建、更新、查询、删除时将被调用，如果任何回调返回错误，
	gorm将停止未来操作并回滚所有更改。

		创建：BeforeSave、BeforeCreate、AfterCreate、AfterSave
		更新：BeforeSave、BeforeUpdate、AfterUpdate、AfterSave
		删除：BeforeDelete、AfterDelete
		查询：AfterFind
*/
// 在 models.go 中使用自定义 callbacks 就不用每个文件重复写这种回调了
// 即使写了也无妨，因为自定义的 callbacks 会在这里后调用，覆盖这的操作
//func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
//	_ = scope.SetColumn("CreatedOn", time.Now().Unix())
//	return nil
//}
//
//func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
//	_ = scope.SetColumn("ModifiedOn", time.Now().Unix())
//	fmt.Println("ModifiedOn=======================>")
//	return nil
//}

//从数据库查询tags
func GetTags(pageNum int, pageSize int, maps interface{}) ([]Tag, error) {
	var (
		tags []Tag
		err  error
	)
	if pageSize > 0 && pageNum > 0 {
		err = db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	} else {
		err = db.Where(maps).Find(&tags).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tags, nil
}

//获取tags总数量
func GetTagTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Tag{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

//查询某个tag是否存在 by name
func ExistTagByName(name string) (bool, error) {
	var tag Tag
	e := db.Select("id").Where("name = ?", name).First(&tag).Error
	if e != nil && e != gorm.ErrRecordNotFound {
		return false, e
	}

	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

//增加tag
func AddTag(name string, state int, createdBy string) error {
	err := db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}).Error
	if err != nil {
		return err
	}

	return nil
}

//查询某个tag是否存在 by id
func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id = ?", id).First(&tag).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

//根据 id 删除某个tag
func DeleteTag(id int) error {
	if err := db.Where("id = ?", id).Delete(&Tag{}).Error; err != nil {
		return err
	}
	return nil
}

//根据 id 修改某个tag信息
func EditTag(id int, data interface{}) error {
	if err := db.Model(&Tag{}).Where("id = ?", id).Update(data).Error; err != nil {
		return err
	}
	return nil
}

// 硬删除标签
func CleanAllTag() (bool, error) {
	// 使用 Unscoped 查找软删除的记录并删除
	if err := db.Unscoped().Where("deleted_on != ?", 0).Delete(&Tag{}).Error; err != nil {
		return false, err
	}
	return true, nil
}
