package scope

import "gorm.io/gorm"

func Paginate(start, limit int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(start).Limit(limit)
	}
}

func PaginateByCol(col string, start any, limit int, desc bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		after := db.Where(col+" > ?", start).Limit(limit)
		if desc {
			return after.Order(col + " DESC")
		}
		return after.Order(col + " ASC")
	}
}
