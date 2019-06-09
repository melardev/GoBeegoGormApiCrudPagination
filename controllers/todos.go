package controllers

import (
	"encoding/json"
	"github.com/melardev/GoBeegoGormApiCrudPagination/dtos"
	"github.com/melardev/GoBeegoGormApiCrudPagination/models"
	"github.com/melardev/GoBeegoGormApiCrudPagination/services"
	"net/http"
	"strconv"
)

type TodosController struct {
	BaseController
}

func (this *TodosController) GetAllTodos() {
	page, pageSize := this.getPagingParams()
	todos, totalTodoCount := services.FetchTodos(page, pageSize)

	this.SendJson(dtos.CreateTodoPagedResponse(this.Ctx.Request.URL.Path, todos, page, pageSize, totalTodoCount))
}

func (this *TodosController) GetAllPendingTodos() {
	page, pageSize := this.getPagingParams()
	todos, totalTodoCount := services.FetchPendingTodos(page, pageSize)
	this.SendJson(dtos.CreateTodoPagedResponse(this.Ctx.Request.URL.Path, todos, page, pageSize, totalTodoCount))
}
func (this *TodosController) GetAllCompletedTodos() {
	page, pageSize := this.getPagingParams()
	todos, totalTodoCount := services.FetchCompletedTodos(page, pageSize)

	this.SendJson(dtos.CreateTodoPagedResponse(this.Ctx.Request.URL.Path, todos, page, pageSize, totalTodoCount))
}

func (this *TodosController) GetTodoById() {
	idStr := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	todo, err := services.FetchById(uint(id))
	if err != nil {
		this.Data["json"] = dtos.CreateErrorDtoWithMessage(err.Error())
		return
	}

	this.SendJson(dtos.GetSuccessTodoDto(&todo))
}

func (this *TodosController) CreateTodo() {
	todoInput := &models.Todo{}

	if err := json.Unmarshal(this.Ctx.Input.RequestBody, todoInput); err != nil {
		// this.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		this.SendJson(dtos.CreateErrorDtoWithMessage(err.Error()))
		return
	}

	todo, err := services.CreateTodo(todoInput.Title, todoInput.Description, todoInput.Completed)
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		this.SendJson(dtos.CreateErrorDtoWithMessage(err.Error()))

		return
	}

	this.SendJson(dtos.CreateTodoCreatedDto(&todo))
}

func (this *TodosController) UpdateTodo() {
	idStr := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		this.SendJson(dtos.CreateErrorDtoWithMessage("You must set an ID"))
		return
	}

	var todoInput models.Todo
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &todoInput); err != nil {
		this.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		this.SendJson(dtos.CreateErrorDtoWithMessage(err.Error()))
		return
	}

	todo, err := services.UpdateTodo(uint(id), todoInput.Title, todoInput.Description, todoInput.Completed)
	if err != nil {
		this.SendJson(dtos.CreateErrorDtoWithMessage(err.Error()))
		return
	}

	this.SendJson(dtos.CreateTodoUpdatedDto(&todo))
}

func (this *TodosController) DeleteTodo() {
	idStr := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		this.SendJson(dtos.CreateErrorDtoWithMessage("You must set an ID"))
		return
	}

	todo, err := services.FetchById(uint(id))

	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		this.SendJson(dtos.CreateErrorDtoWithMessage("todo not found"))
		return
	}

	err = services.DeleteTodo(&todo)

	if err != nil {
		this.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		this.SendJson(dtos.CreateErrorDtoWithMessage("Could not delete Todo"))
		return
	}

	this.SendJson(dtos.CreateSuccessWithMessageDto("todo deleted successfully"))
}

func (this *TodosController) DeleteAllTodos() {
	err := services.DeleteAllTodos()
	if err != nil {
		this.SendJson(dtos.CreateErrorDtoWithMessage(err.Error()))
	}
	this.SendJson(dtos.CreateSuccessWithMessageDto("All Todos deleted succesfully"))
}
