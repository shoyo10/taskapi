package model

type EnumTaskStatus int

const (
	TaskStatusIncomplete EnumTaskStatus = iota
	TaskStatusCompleted
)

type Task struct {
	ID     int            `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name   string         `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Status EnumTaskStatus `json:"status" gorm:"column:status;type:integer;not null"`
}

// TableName ...
func (*Task) TableName() string {
	return "tasks"
}

type TaskFilter struct {
	ID int
}

type UpdateTaskInput struct {
	Name   *string
	Status *EnumTaskStatus
}
