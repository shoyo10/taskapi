package service

import (
	"context"
	"taskapi/internal/model"
	"taskapi/internal/repository/mocks"
	"taskapi/pkg/errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestTaskSvc(t *testing.T) {
	suite.Run(t, new(taskSvcSuite))
}

type taskSvcSuite struct {
	suite.Suite

	mockRepo *mocks.MockRepositorier
	svc      Servicer
}

func (s *taskSvcSuite) SetupSuite() {

}

func (s *taskSvcSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	mockRepo := mocks.NewMockRepositorier(ctrl)
	s.mockRepo = mockRepo
	s.svc = New(mockRepo)
}

func (s *taskSvcSuite) TestCreateTask() {
	ctx := context.Background()
	newTask := model.Task{
		Name:   "test",
		Status: model.TaskStatusIncomplete,
	}
	s.mockRepo.EXPECT().CreateTask(ctx, &newTask).Return(nil).Times(1)
	err := s.svc.CreateTask(ctx, &newTask)
	s.Require().NoError(err)

	s.mockRepo.EXPECT().CreateTask(ctx, gomock.Any()).Return(errors.ErrInternalServerError).Times(1)
	s.Require().ErrorIs(s.svc.CreateTask(ctx, &newTask), errors.ErrInternalServerError)
}

func (s *taskSvcSuite) TestListTask() {
	ctx := context.Background()
	tasks := []model.Task{
		{
			ID:     1,
			Name:   "test",
			Status: model.TaskStatusIncomplete,
		},
	}
	s.mockRepo.EXPECT().ListTask(ctx).Return(tasks, nil).Times(1)
	list, err := s.svc.ListTask(ctx)
	s.Require().NoError(err)
	s.Require().Equal(tasks, list)

	s.mockRepo.EXPECT().ListTask(ctx).Return(nil, errors.ErrInternalServerError).Times(1)
	_, err = s.svc.ListTask(ctx)
	s.Require().ErrorIs(err, errors.ErrInternalServerError)
}

func (s *taskSvcSuite) TestUpdateTask() {
	ctx := context.Background()
	filter := model.TaskFilter{
		ID: 1,
	}
	name := "test"
	status := model.TaskStatusCompleted
	in := model.UpdateTaskInput{
		Name:   &name,
		Status: &status,
	}
	s.mockRepo.EXPECT().UpdateTask(ctx, filter, in).Return(nil).Times(1)
	err := s.svc.UpdateTask(ctx, filter, in)
	s.Require().NoError(err)

	s.mockRepo.EXPECT().UpdateTask(ctx, filter, in).Return(errors.ErrInternalServerError).Times(1)
	s.Require().ErrorIs(s.svc.UpdateTask(ctx, filter, in), errors.ErrInternalServerError)
}

func (s *taskSvcSuite) TestDeleteTask() {
	ctx := context.Background()
	filter := model.TaskFilter{
		ID: 1,
	}
	s.mockRepo.EXPECT().DeleteTask(ctx, filter).Return(nil).Times(1)
	err := s.svc.DeleteTask(ctx, filter)
	s.Require().NoError(err)

	s.mockRepo.EXPECT().DeleteTask(ctx, filter).Return(errors.ErrInternalServerError).Times(1)
	s.Require().ErrorIs(s.svc.DeleteTask(ctx, filter), errors.ErrInternalServerError)
}
