-- name: ListTodoByUserID :many
SELECT 
    todos.*,
    users.name as user_name
FROM todos 
JOIN users ON todos.user_id = users.id
WHERE todos.user_id = $1;

-- name: ListTodos :many
SELECT * FROM todos;
