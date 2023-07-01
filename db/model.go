package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/cryingmouse/data_management_engine/context"
)

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
		for _, attr := range filter.Fields {
			field, ok := reflect.TypeOf(model).FieldByName(attr)
			if ok {
				gormTag := field.Tag.Get("gorm")
				if strings.Contains(gormTag, "column:") {
					attr = strings.Split(gormTag, "column:")[1]
				}

				validAttributes = append(validAttributes, attr)
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

		err = db.Model(&model).Offset((page-1)*pageSize).Limit(pageSize).Find(items, filter.Conditions).Count(&totalCount).Error
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
