package article_service

import (
	"Gin-blog-example/models"
	"Gin-blog-example/pkg/gredis"
	"Gin-blog-example/pkg/logging"
	"Gin-blog-example/service/cache_service"
	"encoding/json"
)

/*
	文章业务提供层
	封装了底层业务，提供向上统一的入口
	分装内容：
		增删改：直接调用 DB 封装模块
		查：使用 DB 查询后，加入 Redis 内存模块，下次读取直接 Redis 内存读取
*/

type ArticleService struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

var cacheTime = 3600

func (ar *ArticleService) Add() error {
	article := map[string]interface{}{
		"tag_id":          ar.TagID,
		"title":           ar.Title,
		"desc":            ar.Desc,
		"content":         ar.Content,
		"created_by":      ar.CreatedBy,
		"cover_image_url": ar.CoverImageUrl,
		"state":           ar.State,
	}
	if err := models.AddArticle(article); err != nil {
		return err
	}
	return nil
}

func (ar *ArticleService) Edit() error {
	return models.EditArticle(ar.ID, map[string]interface{}{
		"tag_id":          ar.TagID,
		"title":           ar.Title,
		"desc":            ar.Desc,
		"content":         ar.Content,
		"cover_image_url": ar.CoverImageUrl,
		"state":           ar.State,
		"modified_by":     ar.ModifiedBy,
	})
}

func (ar *ArticleService) Get() (*models.Article, error) {
	var cacheArticle *models.Article
	cache := cache_service.Article{ID: ar.ID}
	key := cache.GetArticleKey()

	//判断 redis 中是否有，有则直接返回
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			_ = json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	//从数据库查询
	article, err := models.GetArticle(ar.ID)
	if err != nil {
		return nil, err
	}

	//存入 redis 内存缓存，方便下次直接从内存拿取
	_ = gredis.Set(key, article, cacheTime)
	return article, nil
}

func (ar *ArticleService) GetAll() ([]*models.Article, error) {
	var cacheArticles []*models.Article

	cache := cache_service.Article{
		TagID: ar.TagID,
		State: ar.State,

		PageNum:  ar.PageNum,
		PageSize: ar.PageSize,
	}

	key := cache.GetArticlesKey()

	if gredis.Exists(key) {
		data, e := gredis.Get(key)
		if e != nil {
			logging.Info(e)
		} else {
			_ = json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}

	articleList, err := models.GetArticleList(ar.PageNum, ar.PageSize, ar.getMaps())
	if err != nil {
		return nil, err
	}
	_ = gredis.Set(key, articleList, cacheTime)
	return articleList, nil
}

func (ar *ArticleService) Delete() error {
	return models.DeleteArticle(ar.ID)
}

func (ar *ArticleService) ExistByID() (bool, error) {
	return models.ExistArticleByID(ar.ID)
}

func (ar *ArticleService) Count() (int, error) {
	return models.GetArticleTotal(ar.getMaps())
}

// 组装mysql查询条件：未被删除，有状态，有tagid
func (ar *ArticleService) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if ar.State != -1 {
		maps["state"] = ar.State
	}
	if ar.TagID != -1 {
		maps["tag_id"] = ar.TagID
	}

	return maps
}
