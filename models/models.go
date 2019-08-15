package models

import (
	"Gin-blog-example/pkg/setting"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

// 从映射结构体中获取数据库配置信息
func Setup() {
	log.Printf("读取 database 配置项....")

	db, err := gorm.Open(setting.DataBaseSetting.Type,
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			setting.DataBaseSetting.User,
			setting.DataBaseSetting.Password,
			setting.DataBaseSetting.Host,
			setting.DataBaseSetting.Name))

	if err != nil {
		log.Fatalf("models.Setup err:%v", err)
	}

	// 通过定义DefaultTableNameHandler对默认表名应用任何规则,
	// 这里使用自定义前缀+表名
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DataBaseSetting.TablePrefix + defaultTableName
	}

	db.SingularTable(true) //全局禁用表明复数。 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响

	//注册自定义callback
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer db.Close()
}

/*
	通过 GORM 注册自定义 Callbacks 来定制自己的回调驱动，取代 tag.go article.go中的
    BeforeCreate、BeforeUpdate方法，以后类似的文件也不需要写重复的方法
*/

//updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	// 检查是否又含有错误
	if !scope.HasError() {
		nowTime := time.Now().Unix()

		//判断是否包含字段 'CreateOn' 并且该字段是否为空 然后设置值
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

//updateTimeStampForUpdateCallback will set `ModifyTime` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	//scope.Get(...) 根据入参获取设置了字面值的参数，例如本文中是 gorm:update_column ，它会去查找含这个字面值的字段属性
	if _, ok := scope.Get("gorm:update_column"); !ok {
		//设没有指定 update_column 的字段，我们默认在更新回调设置 ModifiedOn 的值.而不是交由GORM这个框架去做
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		// 检查是否手动指定了delete_option
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

		// 获取我们约定的删除字段 'DeletedOn' 如果存在则软删除(逻辑删除)，不存在则直接硬删除
		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),                            //返回应用的表名
				scope.Quote(deletedOnField.DBName),                 //字段在表中的名称
				scope.AddToVars(time.Now().Unix()),                 //添加值作为SQL的参数，也可用于防范SQL注入,这里如果逻辑删除，填入执行删除时候的时间戳
				addExtraSpaceIfExist(scope.CombinedConditionSql()), // 返回组合好的条件SQL
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
