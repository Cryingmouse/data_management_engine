package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/cryingmouse/data_management_engine/common"
)

// parseGormTag 解析 GORM 标签，提取列名
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

// getTagToFieldMap 获取 GORM 标签和字段名的映射关系
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

// Query 根据过滤条件查询数据
func Query(engine *DatabaseEngine, model interface{}, filter *common.QueryFilter, items interface{}) (totalCount int64, err error) {
	db := engine.DB.Debug()

	// 添加关键字条件进行模糊搜索
	for key, value := range filter.Keyword {
		if value != "" {
			db = db.Where(key+" LIKE ?", "%"+value+"%")
		}
	}

	// 根据输入的属性构建动态的 SELECT 语句，或者选择所有属性
	if len(filter.Fields) > 0 {
		// 验证属性是否存在于模型结构体中
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
		db = db.Select(selectStatement)
	} else {
		// Select all attributes if no specific attributes are provided
		db = db.Select("*")
	}

	if filter.Pagination != nil {
		page := filter.Pagination.Page
		pageSize := filter.Pagination.PageSize

		// 获取符合条件的记录总数
		if err = db.Model(model).Where(filter.Conditions).Count(&totalCount).Error; err != nil {
			return totalCount, err
		}

		// 分页查询记录
		err = db.Model(model).Where(filter.Conditions).Offset((page - 1) * pageSize).Limit(pageSize).Find(items).Error
	} else {
		// 不进行分页，查询所有符合条件的记录
		err = db.Model(model).Where(filter.Conditions).Find(items).Error
	}
	return totalCount, err
}

// Delete 根据过滤条件删除数据
func Delete(engine *DatabaseEngine, filter *common.QueryFilter, items interface{}) (err error) {
	db := engine.DB

	// 添加关键字条件进行模糊搜索
	for key, value := range filter.Keyword {
		if value != "" {
			db = db.Where(key+" LIKE ?", "%"+value+"%")
		}
	}

	// 执行删除操作
	if err = db.Unscoped().Where(filter.Conditions).Delete(items).Error; err != nil {
		return fmt.Errorf("failed to delete items:%v in database: %w", items, err)
	}

	return nil
}
