package app

import "time"

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type Task struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"projectId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	AssigneeID  int       `json:"assigneeId"`
	Estimate    int       `json:"estimate"`
	Tags        []string  `json:"tags"`
	DueDate     string    `json:"dueDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Comment struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"taskId"`
	AuthorID  int       `json:"authorId"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

type Database struct {
	Users    []User    `json:"users"`
	Projects []Project `json:"projects"`
	Tasks    []Task    `json:"tasks"`
	Comments []Comment `json:"comments"`
}

type Dashboard struct {
	Users    []User         `json:"users"`
	Projects []Project      `json:"projects"`
	Tasks    []Task         `json:"tasks"`
	Comments []Comment      `json:"comments"`
	Stats    map[string]any `json:"stats"`
	Message  string         `json:"message"`
}
