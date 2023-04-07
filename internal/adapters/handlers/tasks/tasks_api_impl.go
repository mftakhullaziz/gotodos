package tasks

import (
	"github.com/gorilla/mux"
	"gotodo/internal/domain/models/request"
	"gotodo/internal/helpers"
	"gotodo/internal/middleware"
	"gotodo/internal/ports/handlers/api"
	"gotodo/internal/ports/usecases/tasks"
	"net/http"
	"strconv"
	"time"
)

const (
	// formatDatetime is the format string for datetime values.
	formatDatetime = "2006-01-02 15:04:05"
	// messageUserNotAuthorized is the error message for unauthorized user.
	messageUserNotAuthorized = "user account not authorized, please login or sign up!"
	// authHeaderKey is the key for the Authorization header.
	authHeaderKey = "Authorization"
)

type TaskHandlerAPI struct {
	TaskUseCase tasks.TaskUseCase
}

func NewTaskHandlerAPI(taskUseCase tasks.TaskUseCase) api.TaskHandlerAPI {
	return &TaskHandlerAPI{TaskUseCase: taskUseCase}
}

// CreateTaskHandler : do update task based on user authorized
// Params : http.ResponseWriter, *http.Request
func (t TaskHandlerAPI) CreateTaskHandler(writer http.ResponseWriter, requests *http.Request) {
	// Define logger helpers
	log := helpers.LoggerParent()

	// Do get authorization token if any from user login
	token := requests.Header.Get(authHeaderKey)
	authorized, err := middleware.AuthenticateUser(token)
	helpers.LoggerIfError(err)

	// Do check if user account not authorized return empty response
	if authorized == "" {
		responses := helpers.BuildEmptyResponse(messageUserNotAuthorized)
		// Do build write response to response body
		helpers.WriteToResponseBody(writer, &responses)
		return
	}

	// Do convert string authorized to integer
	authorizedUserId, err := strconv.Atoi(authorized)
	helpers.LoggerIfError(err)

	// Do createRequest transform to request body as json
	taskRequest := request.TaskRequest{}
	helpers.ReadFromRequestBody(requests, &taskRequest)
	log.Info("task request body: ", taskRequest)

	// Do get usecase createTask function with param context, updateRequest, userId
	createHandler, err := t.TaskUseCase.CreateTaskUseCase(requests.Context(), taskRequest, authorizedUserId)
	helpers.PanicIfError(err)

	// Do build response handler
	responses := helpers.BuildResponseWithAuthorization(
		createHandler,
		http.StatusCreated,
		int(createHandler.TaskID),
		authorized,
		"create task successful")

	// Do build write response to response body
	helpers.WriteToResponseBody(writer, &responses)
}

// UpdateTaskHandler : do update task based on user authorized and idTask
// Params : http.ResponseWriter, *http.Request
func (t TaskHandlerAPI) UpdateTaskHandler(writer http.ResponseWriter, requests *http.Request) {
	// Define logger helpers
	log := helpers.LoggerParent()

	// Do get authorization token if any from user login
	token := requests.Header.Get(authHeaderKey)
	authorized, err := middleware.AuthenticateUser(token)
	helpers.LoggerIfError(err)

	// Do check if user account not authorized return empty response
	if authorized == "" {
		responses := helpers.BuildEmptyResponse(messageUserNotAuthorized)
		// Do build write response to response body
		helpers.WriteToResponseBody(writer, &responses)
		return
	}

	// Define to get idTask from param
	vars := mux.Vars(requests)
	idTaskVar := vars["task_id"]
	idTask, err := strconv.Atoi(idTaskVar)
	helpers.LoggerIfError(err)

	// Do updateRequest transform to request body as json
	updateRequest := request.TaskRequest{}
	helpers.ReadFromRequestBody(requests, &updateRequest)
	log.Infoln("Update task request: ", updateRequest)

	// Do get usecase updateTask function with param context, updateRequest, idTask
	updateTaskHandler, err := t.TaskUseCase.UpdateTaskUseCase(requests.Context(), updateRequest, idTask)
	helpers.LoggerIfError(err)

	// Do build response handler
	updateTaskResponse := helpers.BuildResponseWithAuthorization(
		updateTaskHandler,
		http.StatusCreated,
		int(updateTaskHandler.TaskID),
		authorized,
		"update task successfully")

	// Do build write response to response body
	helpers.WriteToResponseBody(writer, &updateTaskResponse)
}

func (t TaskHandlerAPI) FindTaskHandlerById(writer http.ResponseWriter, requests *http.Request) {
	// Define logger helpers
	log := helpers.LoggerParent()

	// Do get authorization token if any from user login
	token := requests.Header.Get(authHeaderKey)
	authorized, err := middleware.AuthenticateUser(token)
	helpers.LoggerIfError(err)

	// Do check if user account not authorized return empty response
	if authorized == "" {
		responses := helpers.BuildEmptyResponse(messageUserNotAuthorized)
		// Do build write response to response body
		helpers.WriteToResponseBody(writer, &responses)
		return
	}

	authorizedUserId, err := strconv.Atoi(authorized)
	helpers.LoggerIfError(err)

	// Define to get idTask from param
	vars := mux.Vars(requests)
	idTaskVar := vars["task_id"]
	idTask, err := strconv.Atoi(idTaskVar)
	helpers.LoggerIfError(err)
	log.Infoln("find task by id_task: ", idTask)

	// Do get usecase updateTask function with param context, updateRequest, idTask
	findTaskHandler, err := t.TaskUseCase.FindTaskByIdUseCase(requests.Context(), idTask, authorizedUserId)
	helpers.LoggerIfError(err)
	log.Infoln("find task handler: ", findTaskHandler)

	findTaskHandlerResponse := helpers.BuildResponseWithAuthorization(
		findTaskHandler,
		http.StatusAccepted,
		int(findTaskHandler.TaskID),
		authorized,
		"request find task successful!")

	helpers.WriteToResponseBody(writer, findTaskHandlerResponse)
}

func (t TaskHandlerAPI) FindTaskHandler(writer http.ResponseWriter, requests *http.Request) {
	// Define logger helpers
	log := helpers.LoggerParent()

	// Do get authorization token if any from user login
	token := requests.Header.Get(authHeaderKey)
	authorized, err := middleware.AuthenticateUser(token)
	helpers.LoggerIfError(err)

	// Do check if user account not authorized return empty response
	if authorized == "" {
		responses := helpers.BuildEmptyResponse(messageUserNotAuthorized)
		// Do build write response to response body
		helpers.WriteToResponseBody(writer, &responses)
		return
	}

	authorizedUserId, err := strconv.Atoi(authorized)
	helpers.LoggerIfError(err)

	findAllTaskHandler, err := t.TaskUseCase.FindTaskAllUseCase(requests.Context(), authorizedUserId)
	helpers.LoggerIfError(err)

	tasksSlice := []interface{}{findAllTaskHandler}
	log.Infoln("Tasks: ", tasksSlice)

	findAllTaskHandlerResponse := helpers.BuildAllResponseWithAuthorization(
		tasksSlice[0],
		"request find task successful!",
		len(findAllTaskHandler),
		time.Now().Format(formatDatetime))

	helpers.WriteToResponseBody(writer, findAllTaskHandlerResponse)
}

func (t TaskHandlerAPI) DeleteTaskHandler(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (t TaskHandlerAPI) UpdateTaskStatusHandler(writer http.ResponseWriter, request *http.Request) {
	//TODO implement me
	panic("implement me")
}
