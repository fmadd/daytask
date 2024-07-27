package updateTask

import (
	"daytask/internal/lib/api/response"
	"daytask/internal/lib/logger/sl"
	"daytask/internal/storage"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Description string	 `json:"description"`
	Owner  string `json:"owner"`
	Date   string `json:"date" validate:"required,datetime=2006-01-02"`
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"`
}

type Response struct {
	response.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=TASKUpdater
type TASKUpdater interface {
	UpdateTask(taskID int64, taskName string, taskDescription string, taskOwner string, taskDate string, taskStatus string, taskType string) error
}

// Update task
// @Summary      Update task
// @Description  Update task
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        task   body      Request  true  "updated task"
// @Success      200
// @Failure      400 
// @Failure      404 
// @Failure      500
// @Router       /task [patch]
func New(log *slog.Logger, taskUpdater TASKUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.Update.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req = Request{
			Status: "unstarted",
			Type:   "ordinary",
		}

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		err = taskUpdater.UpdateTask(req.ID, req.Title, req.Description, req.Owner, req.Date, req.Status, req.Type)
		if errors.Is(err, storage.ErrIncorrectDate) {
			log.Info("incorrect date", slog.String("date", req.Date))
			render.JSON(w, r, response.Error("incorrect date"))
			return
		}

		if err != nil {
			log.Error("failed to update task", sl.Err(err))
			render.JSON(w, r, response.Error("failed to update task"))
			return
		}

		log.Info("task update", slog.Int64("id", req.ID))

		render.JSON(w, r, Response{
			Response: response.OK(),
		})
	}
}
