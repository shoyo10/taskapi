package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"taskapi/internal/model"
	"taskapi/internal/service/mocks"
	"taskapi/pkg/echorouter"
	"taskapi/pkg/errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

func TestTaskHanlder(t *testing.T) {
	suite.Run(t, new(taskHandlerSuite))
}

type taskHandlerSuite struct {
	suite.Suite

	mockSvc *mocks.MockServicer
	handler Handler
	app     *fx.App
	e       *echo.Echo
}

func (s *taskHandlerSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.mockSvc = mocks.NewMockServicer(ctrl)
	s.handler = NewHandler(s.mockSvc)
	s.app = fx.New(
		fx.Supply(&echorouter.Config{}),
		fx.Provide(echorouter.FxNewEcho),
		fx.Populate(&s.e),
	)
	SetRoutes(s.e, s.handler)
	err := s.app.Start(context.Background())
	s.Require().NoError(err)
}

func (s *taskHandlerSuite) TearDownSuite() {
	s.app.Stop(context.Background())
}

func (s *taskHandlerSuite) TestListTask() {
	s.mockSvc.EXPECT().ListTask(gomock.Any()).Return(nil, nil).Times(1)

	rec := request(http.MethodGet, "/tasks", nil, s.e)
	s.Require().Equal(http.StatusOK, rec.Code)
	var resp1 listTaskResp
	json.Unmarshal(rec.Body.Bytes(), &resp1)
	s.Require().Len(resp1.Data, 0)

	tasks := []model.Task{
		{
			ID:     1,
			Name:   "task1",
			Status: model.TaskStatusIncomplete,
		},
		{
			ID:     1,
			Name:   "task2",
			Status: model.TaskStatusCompleted,
		},
	}
	s.mockSvc.EXPECT().ListTask(gomock.Any()).Return(tasks, nil).Times(1)
	rec = request(http.MethodGet, "/tasks", nil, s.e)
	var resp2 listTaskResp
	json.Unmarshal(rec.Body.Bytes(), &resp2)
	s.Require().Len(resp2.Data, 2)
}

func (s *taskHandlerSuite) TestCreateTask() {
	type validReqTest struct {
		req            createTaskReq
		expectHTTPCode int
	}

	reqValidCase := []validReqTest{
		{
			req: createTaskReq{
				Name:   "task1",
				Status: 0,
			},
			expectHTTPCode: http.StatusOK,
		},
		{
			req: createTaskReq{
				Name:   "task2",
				Status: 1,
			},
			expectHTTPCode: http.StatusOK,
		},
	}

	for _, c := range reqValidCase {
		expectNewTask := &model.Task{
			Name:   c.req.Name,
			Status: model.EnumTaskStatus(c.req.Status),
		}
		s.mockSvc.EXPECT().CreateTask(gomock.Any(), expectNewTask).DoAndReturn(func(_ context.Context, in *model.Task) error {
			in.ID = 1
			return nil
		}).Times(1)
		reqBody, _ := json.Marshal(c.req)
		rec := request(http.MethodPost, "/tasks", bytes.NewReader(reqBody), s.e)
		s.Require().Equal(c.expectHTTPCode, rec.Code)
		var resp createTaskResp
		json.Unmarshal(rec.Body.Bytes(), &resp)
		s.Require().Equal(1, resp.Data.ID)
	}

	type invalidReqTest struct {
		description    string
		req            createTaskReq
		expectHTTPCode int
		expectErrCode  string
	}
	reqInvalidCase := []invalidReqTest{
		{
			description: "name should not be empty",
			req: createTaskReq{
				Name: "",
			},
			expectHTTPCode: http.StatusBadRequest,
			expectErrCode:  errors.ErrInvalidInput.Code,
		},
		{
			description: "name legth should be less than 32",
			req: createTaskReq{
				Name: "qwertyuiopasdfghjklzxcvbnmqwertyu",
			},
			expectHTTPCode: http.StatusBadRequest,
			expectErrCode:  errors.ErrInvalidInput.Code,
		},
		{
			description: "status should be 0 or 1",
			req: createTaskReq{
				Name:   "task1",
				Status: 2,
			},
			expectHTTPCode: http.StatusBadRequest,
			expectErrCode:  errors.ErrInvalidInput.Code,
		},
	}
	for _, c := range reqInvalidCase {
		reqBody, _ := json.Marshal(c.req)
		rec := request(http.MethodPost, "/tasks", bytes.NewReader(reqBody), s.e)
		s.Require().Equal(c.expectHTTPCode, rec.Code)
		var httpErr errors.HTTPError
		_ = json.Unmarshal(rec.Body.Bytes(), &httpErr)
		s.Require().Equal(c.expectErrCode, httpErr.Code)
	}
}

func (s *taskHandlerSuite) TestUpdateTask() {
	type invalidReqTest struct {
		description    string
		req            updateTaskReq
		path           string
		expectHTTPCode int
		expectErrCode  string
	}
	reqInvalidCase := []invalidReqTest{
		{
			description: "name should not be empty",
			req: updateTaskReq{
				Name:   "",
				Status: 0,
			},
			path:           "/tasks/1",
			expectHTTPCode: http.StatusBadRequest,
			expectErrCode:  errors.ErrInvalidInput.Code,
		},
		{
			description: "name legth should be less than 32",
			req: updateTaskReq{
				Name:   "qwertyuiopasdfghjklzxcvbnmqwertyu",
				Status: 0,
			},
			path:           "/tasks/1",
			expectHTTPCode: http.StatusBadRequest,
			expectErrCode:  errors.ErrInvalidInput.Code,
		},
		{
			description: "status should be 0 or 1",
			req: updateTaskReq{
				Name:   "task1",
				Status: 2,
			},
			path:           "/tasks/1",
			expectHTTPCode: http.StatusBadRequest,
			expectErrCode:  errors.ErrInvalidInput.Code,
		},
		{
			description: "path id out of range",
			req: updateTaskReq{
				Name:   "task1",
				Status: 0,
			},
			path:           "/tasks/12345678912345123456789067890",
			expectHTTPCode: http.StatusBadRequest,
			expectErrCode:  errors.ErrInvalidInput.Code,
		},
	}
	for _, c := range reqInvalidCase {
		reqBody, _ := json.Marshal(c.req)
		rec := request(http.MethodPut, c.path, bytes.NewReader(reqBody), s.e)
		s.Require().Equal(c.expectHTTPCode, rec.Code)
		var httpErr errors.HTTPError
		_ = json.Unmarshal(rec.Body.Bytes(), &httpErr)
		s.Require().Equal(c.expectErrCode, httpErr.Code)
	}

	// test success case
	req := updateTaskReq{
		Name:   "task1",
		Status: 1,
	}
	reqBody, _ := json.Marshal(req)
	expectUpdateReq := model.UpdateTaskInput{
		Name:   &req.Name,
		Status: (*model.EnumTaskStatus)(&req.Status),
	}
	expectFilter := model.TaskFilter{
		ID: 1,
	}
	s.mockSvc.EXPECT().UpdateTask(gomock.Any(), expectFilter, expectUpdateReq).Return(nil).Times(1)
	rec := request(http.MethodPut, "/tasks/1", bytes.NewReader(reqBody), s.e)
	s.Require().Equal(http.StatusOK, rec.Code)
}

func (s *taskHandlerSuite) TestDeleteTask() {
	type invalidReqTest struct {
		description    string
		path           string
		expectHTTPCode int
		expectErrCode  string
	}
	reqInvalidCase := []invalidReqTest{
		{
			description:    "path id out of range",
			path:           "/tasks/12345678912345123456789067890",
			expectHTTPCode: http.StatusBadRequest,
			expectErrCode:  errors.ErrInvalidInput.Code,
		},
	}
	for _, c := range reqInvalidCase {
		rec := request(http.MethodDelete, c.path, nil, s.e)
		s.Require().Equal(c.expectHTTPCode, rec.Code)
		var httpErr errors.HTTPError
		_ = json.Unmarshal(rec.Body.Bytes(), &httpErr)
		s.Require().Equal(c.expectErrCode, httpErr.Code)
	}

	// test success case
	s.mockSvc.EXPECT().DeleteTask(gomock.Any(), model.TaskFilter{ID: 1}).Return(nil).Times(1)
	rec := request(http.MethodDelete, "/tasks/1", nil, s.e)
	s.Require().Equal(http.StatusOK, rec.Code)

}

func request(method, path string, body io.Reader, e *echo.Echo) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}
