package delete

import (
	"daytask/internal/lib/api/response"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"daytask/internal/lib/logger/sl"
	_ "daytask/internal/storage"
	_ "errors"

)

type Request struct{
	ID       int64	 `json:"id"` 
}

type Response struct{
	response.Response	
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=TASKDeleter
type TASKDeleter interface{
	DeleteTask(id int64) (error)
}

// Delete task
// @Summary      Delete task
// @Description  Delete task
// @Tags         task
// @Accept       json
// @Produce      json
// @Param        id   body      int  true  "task ID"
// @Success      200  
// @Failure      400  
// @Failure      404  
// @Failure      500  
// @Router       /task [delete]
func New(log *slog.Logger, taskDeleter TASKDeleter) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		const op = "handlers.task.delete"

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
		
		if err := validator.New().Struct(req); err != nil{
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		err = taskDeleter.DeleteTask(req.ID)
		
		if err != nil {
			log.Error("failed to delete task", sl.Err(err))
			render.JSON(w, r, response.Error("failed to delete task"))
			return
		}

		log.Info("task deleted")

		render.JSON(w, r, Response{ 
			Response: response.OK(),
		})
	}
}