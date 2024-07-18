package sqlite

import (
	"context"
	"log"
	"taskapi/internal/model"
	"taskapi/internal/repository"
	"taskapi/pkg/errors"
	"taskapi/pkg/sqlite"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestTaskDB(t *testing.T) {
	suite.Run(t, new(taskDBSuite))
}

type taskDBSuite struct {
	suite.Suite
	conn *gorm.DB
	repo repository.Repositorier
}

func (s *taskDBSuite) SetupSuite() {
	conn, err := sqlite.New()
	s.Require().NoError(err)
	s.conn = conn
	repo, err := New(conn)
	s.Require().NoError(err)
	s.repo = repo
}

func (s *taskDBSuite) SetupTest() {
	if err := s.conn.Delete(&model.Task{}, "1 = 1").Error; err != nil {
		log.Println("SetupTest delete table failed", err.Error())
	}
}

func (s *taskDBSuite) TearDownSuite() {
	if err := s.conn.Delete(&model.Task{}, "1 = 1").Error; err != nil {
		log.Println("TearDownSuite delete table failed", err.Error())
	}
}

func (s *taskDBSuite) TestCreateTask() {
	newTask := model.Task{
		Name:   "test",
		Status: model.TaskStatusIncomplete,
	}
	ctx := context.Background()
	err := s.repo.CreateTask(ctx, &newTask)
	s.Require().NoError(err)
	s.Require().Equal(1, newTask.ID)

	// Check if the task is created
	var dbTask model.Task
	err = s.conn.First(&dbTask, newTask.ID).Error
	s.Require().NoError(err)
	s.Require().Equal(newTask.Name, dbTask.Name)
	s.Require().Equal(newTask.Status, dbTask.Status)

	// Check if omit id value when creat task
	newTask = model.Task{
		ID:     100,
		Name:   "id omit",
		Status: model.TaskStatusCompleted,
	}
	err = s.repo.CreateTask(ctx, &newTask)
	s.Require().NoError(err)
	s.Require().Equal(2, newTask.ID)

	var dbTask2 model.Task
	err = s.conn.First(&dbTask2, newTask.ID).Error
	s.Require().NoError(err)
	s.Require().Equal(newTask.Name, dbTask2.Name)
	s.Require().Equal(model.TaskStatusCompleted, dbTask2.Status)
}

func (s *taskDBSuite) TestListTask() {
	ctx := context.Background()
	resp, err := s.repo.ListTask(ctx)
	s.Require().NoError(err)
	s.Require().Empty(resp)

	task1 := model.Task{Name: "task1", Status: model.TaskStatusIncomplete}
	s.Require().NoError(s.conn.Create(&task1).Error)

	task2 := model.Task{Name: "task2", Status: model.TaskStatusCompleted}
	s.Require().NoError(s.conn.Create(&task2).Error)

	resp, err = s.repo.ListTask(ctx)
	s.Require().NoError(err)
	s.Require().Len(resp, 2)
}

func (s *taskDBSuite) TestUpdateTask() {
	ctx := context.Background()
	err := s.repo.UpdateTask(ctx, repository.TaskFilter{ID: 1}, repository.UpdateTaskInput{})
	s.Require().Error(err)
	s.Require().ErrorIs(err, errors.ErrResourceNotFound)

	task1 := model.Task{Name: "task1", Status: model.TaskStatusCompleted}
	s.Require().NoError(s.conn.Create(&task1).Error)
	var dbTask model.Task
	s.Require().NoError(s.conn.First(&dbTask, task1.ID).Error)
	s.Require().Equal(task1.Name, dbTask.Name)
	s.Require().Equal(task1.Status, dbTask.Status)

	name := "task1 updated"
	status := model.TaskStatusIncomplete
	err = s.repo.UpdateTask(ctx, repository.TaskFilter{ID: task1.ID}, repository.UpdateTaskInput{
		Name:   &name,
		Status: &status,
	})
	s.Require().NoError(err)

	var updatedDbTask model.Task
	s.Require().NoError(s.conn.First(&updatedDbTask, task1.ID).Error)
	s.Require().Equal(name, updatedDbTask.Name)
	s.Require().Equal(status, updatedDbTask.Status)
}

func (s *taskDBSuite) TestDeleteTask() {
	ctx := context.Background()

	task1 := model.Task{Name: "task1", Status: model.TaskStatusCompleted}
	s.Require().NoError(s.conn.Create(&task1).Error)
	var dbTask model.Task
	s.Require().NoError(s.conn.First(&dbTask, task1.ID).Error)
	s.Require().Equal(task1.Name, dbTask.Name)
	s.Require().Equal(task1.Status, dbTask.Status)

	err := s.repo.DeleteTask(ctx, repository.TaskFilter{ID: task1.ID})
	s.Require().NoError(err)

	err = s.conn.First(&dbTask, task1.ID).Error
	s.Require().ErrorIs(err, gorm.ErrRecordNotFound)
}
