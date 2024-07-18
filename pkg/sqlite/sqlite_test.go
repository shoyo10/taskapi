package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	db, err := New()
	assert.Nil(t, err)

	err = db.AutoMigrate(&ABC{})
	assert.Nil(t, err)
	data := ABC{
		Name: "test",
	}
	err = db.Create(&data).Error
	assert.Nil(t, err)
	assert.Equal(t, 1, data.ID)
}

type ABC struct {
	ID   int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"column:name;type:varchar(255);not null"`
}

// TableName ...
func (*ABC) TableName() string {
	return "abc"
}
