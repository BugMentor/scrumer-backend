package graph

import "gorm.io/gorm"

// GormDB is an interface that defines the GORM DB methods used by our resolvers.
// This allows us to mock the DB for testing purposes.
type GormDB interface {
	Create(value interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Find(dest interface{}, conds ...interface{}) (tx *gorm.DB) // Added Find method
	Save(value interface{}) (tx *gorm.DB)
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
	Preload(query string, args ...interface{}) (tx *gorm.DB)
	Model(value interface{}) (tx *gorm.DB)
	Association(column string) *gorm.Association
}

// Ensure that *gorm.DB implements the GormDB interface
var _ GormDB = (*gorm.DB)(nil)
