package tag_service

import (
	"Gin-blog-example/models"
	"Gin-blog-example/pkg/export"
	"Gin-blog-example/pkg/file"
	"Gin-blog-example/pkg/gredis"
	"Gin-blog-example/pkg/logging"
	"Gin-blog-example/pkg/setting"
	"Gin-blog-example/service/cache_service"
	"encoding/json"
	"github.com/Unknwon/com"
	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize"
	"io"
	"strconv"
	"time"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}

	return models.EditTag(t.ID, data)
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)

	cache := cache_service.Tag{
		State: t.State,

		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := cache.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			_ = json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	_ = gredis.Set(key, tags, 3600)
	return tags, nil
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}

// Export 导出标签 Excel 文件
func (t *Tag) Export() (string, error) {

	tags, e := t.GetAll()
	if e != nil {
		return "", e
	}

	tagFile := xlsx.NewFile()
	//新增一个 sheet
	sheet, e := tagFile.AddSheet("标签信息")
	if e != nil {
		return "", e
	}

	//定义单元格名称
	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	//新增一行
	row := sheet.AddRow()

	var cell *xlsx.Cell

	//将单元格名称插入到第一行
	for _, title := range titles {
		//将单元格出入当前行
		cell = row.AddCell()
		cell.Value = title
	}

	for _, v := range tags {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			com.Date(int64(v.CreatedOn), setting.AppSetting.TimeFormat),
			v.ModifiedBy,
			com.Date(int64(v.ModifiedOn), setting.AppSetting.TimeFormat),
		}

		row = sheet.AddRow()

		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	time := strconv.Itoa(int(time.Now().Unix()))

	filename := "tags-" + time + export.EXT
	dirFullPath := export.GetExcelFullPath()
	e = file.IsNotExistMkDir(dirFullPath)
	if e != nil {
		return "", e
	}

	fullPath := dirFullPath + filename
	e = tagFile.Save(fullPath)
	if e != nil {
		return "", e
	}

	return filename, nil
}

//导入标签
func (t *Tag) Import(r io.Reader) error {
	xlsx, e := excelize.OpenReader(r)
	if e != nil {
		return e
	}

	//获取某张 sheet 的所有行信息
	rows, e := xlsx.GetRows("标签信息")
	if e != nil {
		return e
	}

	for irow, row := range rows {
		if irow > 0 {
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}

			models.AddTag(data[1], 1, data[2])
		}
	}

	return nil
}
