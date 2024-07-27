package getAllTasks

import (
	"daytask/internal/lib/api/response"
	"daytask/internal/lib/logger/sl"
	"daytask/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct{
	Owner		string	 `json:"owner"`
}

type Response struct{
	response.Response
	Quantity	int		 		 `json:"quantity"`
	Tasks       []storage.Task	 `json:"tasks"`		
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=TASKGetterAll
type TASKGetterAll interface {
	GetAllTasks(taskOwner string) ([]storage.Task, error)
}
// Get all tasks
// @Summary      Get all tasks
// @Description  Gives tasks for the whole time
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        owner   body      string  true  "owner's login"
// @Success      200  {object}	Response "Quantity and Tasks array" 
// @Failure      400  
// @Failure      404  
// @Failure      500  
// @Router       /task/all [get]
func New(log *slog.Logger, taskGetterAll TASKGetterAll) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		const op = "handlers.task.GetAllTasks.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil{
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}
		
		log.Info("request body decoded", slog.Any("request", req))
		
		tasks, err := taskGetterAll.GetAllTasks(req.Owner)


		if err != nil {
			log.Error("failed to get all tasks", sl.Err(err))
			render.JSON(w, r, response.Error("failed to get all tasks"))
			return
		}

		log.Info("all task get", slog.Any("quantity", len(tasks)))

		render.JSON(w, r, Response{ 
			Response: response.OK(),
			Quantity: len(tasks),
			Tasks: tasks,
		})
	}
}