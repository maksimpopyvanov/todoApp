package repository

import (
	"fmt"
	"todo"

	"github.com/jmoiron/sqlx"
)

type TodoItemRepository struct {
	db *sqlx.DB
}

func NewTodoItemRepository(db *sqlx.DB) *TodoItemRepository {
	return &TodoItemRepository{
		db: db,
	}
}

func (r *TodoItemRepository) Create(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		return 0, err
	}

	createListsItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listItemsTable)
	_, err = tx.Exec(createListsItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()

}

func (r *TodoItemRepository) GetAll(listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON ti.id = li.item_id WHERE li.list_id = $1",
		todoItemsTable, listItemsTable)

	err := r.db.Select(&items, query, listId)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemRepository) GetById(userId int, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON ti.id = li.item_id
	INNER JOIN %s ul ON ul.list_id = li.list_id WHERE ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, listItemsTable, usersListsTable)

	err := r.db.Get(&item, query, itemId, userId)
	if err != nil {
		return item, err
	}

	return item, nil
}
