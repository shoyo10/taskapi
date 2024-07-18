package http

import (
	"strings"
	"taskapi/internal/model"
	"taskapi/internal/service"
	"taskapi/pkg/errors"

	"net/http"

	"github.com/labstack/echo/v4"
)

type handler struct {
	svc service.Servicer
}

// NewHandler create Handler instance
func NewHandler(svc service.Servicer) Handler {
	return &handler{
		svc: svc,
	}
}

type listTaskResp struct {
	Data []listTaskRespData `json:"data"`
}

type listTaskRespData struct {
	ID     int                  `json:"id"`
	Name   string               `json:"name"`
	Status model.EnumTaskStatus `json:"status"`
}

// @Title  ListTask
// @Description list all tasks
// @Success 200 {object} listTaskResp
// @Router /tasks [get]
func (h *handler) ListTask(c echo.Context) error {
	ctx := c.Request().Context()
	result, err := h.svc.ListTask(ctx)
	if err != nil {
		return err
	}
	resp := listTaskResp{
		Data: make([]listTaskRespData, len(result)),
	}
	for i, r := range result {
		resp.Data[i] = listTaskRespData{
			ID:     r.ID,
			Name:   r.Name,
			Status: r.Status,
		}
	}
	return c.JSON(http.StatusOK, resp)
}

type createTaskReq struct {
	Name   string `json:"name" validate:"gte=1,lte=32"`
	Status int    `json:"status" validate:"oneof=0 1"`
}

type createTaskResp struct {
	Data createTaskRespData `json:"data"`
}

type createTaskRespData struct {
	ID int `json:"id"`
}

// @Title  CreateTask
// @Description create a task
// @Param reqBody body createTaskReq true "task fields"
// @Success 200 {object} createTaskResp "task id"
// @Failure 400 object errors.HTTPError
// @Router /tasks [post]
func (h *handler) CreateTask(c echo.Context) error {
	ctx := c.Request().Context()

	var req createTaskReq
	err := c.Bind(&req)
	if err != nil {
		return errors.Wrap(errors.ErrInvalidInput, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return errors.Wrap(errors.ErrInvalidInput, err.Error())
	}

	req.Name = strings.TrimSpace(req.Name)

	newTask := &model.Task{
		Name:   req.Name,
		Status: model.EnumTaskStatus(req.Status),
	}
	err = h.svc.CreateTask(ctx, newTask)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, createTaskResp{Data: createTaskRespData{ID: newTask.ID}})
}

type updateTaskReq struct {
	Name   string `json:"name" validate:"gte=1,lte=32"`
	Status int    `json:"status" validate:"oneof=0 1"`
}

type updateTaskParam struct {
	ID int `param:"id" validate:"required"`
}

// @Title  UpdateTask
// @Description update a task
// @Param id path int true "task id"
// @Param reqBody body updateTaskReq true "update task fields"
// @Success 200
// @Failure 400 object errors.HTTPError
// @Failure 404 object errors.HTTPError
// @Router /tasks/{id} [put]
func (h *handler) UpdateTask(c echo.Context) error {
	ctx := c.Request().Context()

	var req updateTaskReq
	err := c.Bind(&req)
	if err != nil {
		return errors.Wrap(errors.ErrInvalidInput, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return errors.Wrap(errors.ErrInvalidInput, err.Error())
	}

	var pathParam updateTaskParam
	err = (&echo.DefaultBinder{}).BindPathParams(c, &pathParam)
	if err != nil {
		return errors.Wrap(errors.ErrInvalidInput, err.Error())
	}
	if err := c.Validate(&pathParam); err != nil {
		return errors.Wrap(errors.ErrInvalidInput, err.Error())
	}

	updateReq := model.UpdateTaskInput{
		Name:   &req.Name,
		Status: (*model.EnumTaskStatus)(&req.Status),
	}
	err = h.svc.UpdateTask(ctx, model.TaskFilter{ID: pathParam.ID}, updateReq)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

type deleteTaskParam struct {
	ID int `param:"id" validate:"required"`
}

// @Title  DeleteTask
// @Description delete a task
// @Param id path int true "task id"
// @Success 200
// @Failure 400 object errors.HTTPError
// @Failure 404 object errors.HTTPError
// @Router /tasks/{id} [delete]
func (h *handler) DeleteTask(c echo.Context) error {
	ctx := c.Request().Context()

	var pathParam deleteTaskParam
	err := (&echo.DefaultBinder{}).BindPathParams(c, &pathParam)
	if err != nil {
		return errors.Wrap(errors.ErrInvalidInput, err.Error())
	}
	if err := c.Validate(&pathParam); err != nil {
		return errors.Wrap(errors.ErrInvalidInput, err.Error())
	}

	err = h.svc.DeleteTask(ctx, model.TaskFilter{ID: pathParam.ID})
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
