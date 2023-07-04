package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/cryingmouse/data_management_engine/context"
)

func parseGormTag(tag string) string {
	if tag == "" {
		return ""
	}

	tags := strings.Split(tag, ";")
	for _, t := range tags {
		if strings.HasPrefix(t, "column:") {
			column := strings.TrimPrefix(t, "column:")
			return column
		}
	}

	return ""
}

func getTagToFieldMap(model interface{}) map[string]string {
	tagToFieldMap := make(map[string]string)

	modelType := reflect.TypeOf(model)
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		gormTag := field.Tag.Get("gorm")

		// 解析 GORM 标签
		columnName := parseGormTag(gormTag)

		// 将 GORM 标签和属性名映射到 map
		if columnName != "" {
			tagToFieldMap[columnName] = field.Name
		}
	}

	return tagToFieldMap
}

func Query(engine *DatabaseEngine, model interface{}, filter *context.QueryFilter, items interface{}) (totalCount int64, err error) {
	db := engine.DB.Debug()

	// Add the keyword to the conditions for the fuzzy search
	for key, value := range filter.Keyword {
		if value != "" {
			db = db.Where(key+" LIKE ?", "%"+value+"%")
		}
	}

	// Build the SELECT statement dynamically based on the input attributes or retrieve all attributes
	if len(filter.Fields) > 0 {
		// Validate attributes exist in the Directory struct
		var validAttributes []string
		tagToFieldMap := getTagToFieldMap(model)
		for _, attr := range filter.Fields {
			if _, ok := tagToFieldMap[attr]; ok {
				validAttributes = append(validAttributes, attr)
			} else {
				return totalCount, fmt.Errorf("invalid attribute:%v found", attr)
			}
		}
		if len(validAttributes) == 0 {
			return totalCount, errors.New("no valid attributes found")
		}
		// Use the provided attributes
		selectStatement := strings.Join(validAttributes, ", ")
		// selectStatement := strings.Join(filter.attributes, ", ")
		db = db.Select(selectStatement)
	} else {
		// Select all attributes if no specific attributes are provided
		db = db.Select("*")
	}

	if filter.Pagination != nil {
		page := filter.Pagination.Page
		pageSize := filter.Pagination.PageSize

		if err = db.Model(&model).Find(items, filter.Conditions).Count(&totalCount).Error; err != nil {
			return totalCount, err
		}

		err = db.Model(&model).Offset((page-1)*pageSize).Limit(pageSize).Find(items, filter.Conditions).Error
	} else {
		err = db.Model(&model).Find(items, filter.Conditions).Error
	}
	return totalCount, err
}

func Delete(engine *DatabaseEngine, filter *context.QueryFilter, items interface{}) (err error) {
	db := engine.DB.Debug().Where("1 = 1")

	for key, value := range filter.Keyword {
		if value != "" {
			db = db.Where(key+" LIKE ?", "%"+value+"%")
		}
	}

	if err = db.Unscoped().Delete(items, filter.Conditions).Error; err != nil {
		return fmt.Errorf("failed to delete items:%v in database: %w", items, err)
	}

	return
}
