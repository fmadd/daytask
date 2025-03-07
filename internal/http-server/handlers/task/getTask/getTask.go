package getTask

import (
	"daytask/internal/lib/api/response"
	"log/slog"
	"net/http"

	"daytask/internal/lib/logger/sl"
	"daytask/internal/storage"
	"errors"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Owner string `json:"owner"`
	Date  string `json:"date" validate:"required,datetime=2006-01-02"`
}

type Response struct {
	response.Response
	Quantity int            `json:"quantity,omitempty"`
	Tasks    []storage.Task `json:"tasks,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=TASKGetter
type TASKGetter interface {
	GetTaskForDay(taskOwner string, taskDate string) ([]storage.Task, error)
}

// Get the day's tasks
// @Summary      Get tasks
// @Description  Get the day's tasks
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        owner   body      string  true  "owner's login"
// @Param        date   body      string  true  "date"
// @Success      200  {object} Response "Quantity and Tasks array"
// @Failure      400 
// @Failure      404 
// @Failure      500 
// @Router       /task/day [get]
func New(log *slog.Logger, taskGetter TASKGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.gettask"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

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

		tasks, err := taskGetter.GetTaskForDay(req.Owner, req.Date)
		if errors.Is(err, storage.ErrIncorrectDate) {
			log.Info("incorrect date", slog.String("date", req.Date))
			render.JSON(w, r, response.Error("incorrect date"))
			return
		}

		if err != nil {
			log.Error("failed to get task", sl.Err(err))
			render.JSON(w, r, response.Error("failed to get task"))
			return
		}

		log.Info("task get", slog.Any("quantity", len(tasks)))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Quantity: len(tasks),
			Tasks:    tasks,
		})
	}
}
