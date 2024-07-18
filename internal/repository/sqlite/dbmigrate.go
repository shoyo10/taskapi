package sqlite

import "taskapi/internal/model"

func (r *repo) Migrate() error {
	return r.conn.AutoMigrate(&model.Task{})
}
