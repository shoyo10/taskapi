package sqlite

import (
	"context"
	"taskapi/internal/model"
	"taskapi/internal/repository"
	"taskapi/pkg/errors"

	"gorm.io/gorm"
)

type repo struct {
	conn *gorm.DB
}

// New ...
func New(conn *gorm.DB) (repository.Repositorier, error) {
	r := &repo{
		conn: conn,
	}
	if err := r.Migrate(); err != nil {
		return nil, errors.WithStack(err)
	}
	return r, nil
}

func (r *repo) Ctx(ctx context.Context) *gorm.DB {
	return r.conn.WithContext(ctx)

}

func (r *repo) CreateTask(ctx context.Context, task *model.Task) error {
	err := r.Ctx(ctx).Omit("id").Create(task).Error
	if err != nil {
		return errors.Wrapf(errors.ErrInternalServerError, "%v", err)
	}
	return nil
}

func (r *repo) ListTask(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	err := r.Ctx(ctx).Find(&tasks).Error
	if err != nil {
		return nil, errors.Wrapf(errors.ErrInternalServerError, "%v", err)
	}
	return tasks, nil
}

func (r *repo) UpdateTask(ctx context.Context, filter model.TaskFilter, in model.UpdateTaskInput) error {
	result := r.Ctx(ctx).Model(&model.Task{}).Where("id = ?", filter.ID).Omit("id").Updates(&in)
	rawEffected, err := result.RowsAffected, result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrapf(errors.ErrResourceNotFound, "%v", err)
		}
		return errors.Wrapf(errors.ErrInternalServerError, "%v", err)
	}
	if rawEffected == 0 {
		return errors.Wrapf(errors.ErrResourceNotFound, "task with filter %v not found", filter)

	}
	return nil
}

func (r *repo) DeleteTask(ctx context.Context, filter model.TaskFilter) error {
	err := r.Ctx(ctx).Where("id = ?", filter.ID).Delete(&model.Task{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrapf(errors.ErrResourceNotFound, "%v", err)
		}
		return errors.Wrapf(errors.ErrInternalServerError, "%v", err)
	}
	return nil
}
