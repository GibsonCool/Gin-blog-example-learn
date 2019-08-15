package models

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
func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
	return
}

//获取tags总数量
func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)
	return
}

//查询某个tag是否存在 by name
func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name = ?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

//增加tag
func AddTag(name string, state int, createdBy string) bool {
	db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	})

	return true
}

//查询某个tag是否存在 by id
func ExistTagByID(id int) bool {
	var tag Tag
	db.Select("id").Where("id = ?", id).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

//根据 id 删除某个tag
func DeleteTag(id int) bool {
	db.Where("id = ?", id).Delete(&Tag{})
	return true
}

//根据 id 修改某个tag信息
func EditTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id = ?", id).Update(data)
	return true
}

// 硬删除标签
func CleanAllTag() bool {
	// 使用 Unscoped 查找软删除的记录并删除
	db.Unscoped().Where("deleted_on != ?", 0).Delete(&Tag{})
	return true
}
