package database

import "github.com/iamrahultanwar/heeico/model"

func AliasExist(urlShorter model.URLShortener) bool {
	var count int64
	DB.Model(urlShorter).Where("alias", urlShorter.Alias).Count(&count)
	return count > 0
}
