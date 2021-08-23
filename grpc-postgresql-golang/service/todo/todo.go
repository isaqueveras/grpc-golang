package todo

import (
	"context"

	"github.com/go-pg/pg"
)

type Service struct {
	DB *pg.DB
}

func (pg Service) CreateTodo(ctx context.Context, req *todo.CreateTodoRequest) {

}
