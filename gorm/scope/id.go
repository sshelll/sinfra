package scope

import "gorm.io/gorm"

func SelectAndOrderById(id any) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where("id = ?", id).Order("id")
	}
}

func SelectAndOrderByCol(col string, val any) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Where(col+" = ?", val).Order(col)
	}
}
