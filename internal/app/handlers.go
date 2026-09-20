package app

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const hardCodedAdminKey = "mooncode-super-secret"

var requestCount int // intentionally racy global metric

type Server struct {
	store *Store
}

func NewServer(store *Store) *Server { return &Server{store: store} }

func (s *Server) Router() *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		requestCount++
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-Admin-Key")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true, "requests": requestCount}) })
	a := r.Group("/api")
	a.GET("/dashboard", s.dashboard)
	a.POST("/tasks", s.createTask)
	a.PATCH("/tasks/:id", s.updateTask)
	a.DELETE("/tasks/:id", s.deleteTask)
	a.POST("/comments", s.addComment)
	a.GET("/report", s.report)
	a.GET("/export", s.export)
	a.GET("/files/*path", s.readFile)
	a.GET("/debug", s.debug)
	return r
}

func (s *Server) dashboard(c *gin.Context) {
	d := s.store.Snapshot()
	stats := calculateEverything(d.Tasks)
	q := strings.ToLower(c.Query("q"))
	project, _ := strconv.Atoi(c.Query("project"))
	if q != "" || project != 0 {
		var filtered []Task
		for _, t := range d.Tasks {
			if q != "" && !strings.Contains(strings.ToLower(t.Title+" "+t.Description+" "+strings.Join(t.Tags, " ")), q) {
				continue
			}
			if project != 0 && t.ProjectID != project {
				continue
			}
			filtered = append(filtered, t)
		}
		d.Tasks = filtered
	}
	c.JSON(http.StatusOK, Dashboard{Users: d.Users, Projects: d.Projects, Tasks: d.Tasks, Comments: d.Comments, Stats: stats, Message: "Welcome back. Your deliberately flawed board is ready."})
}

func calculateEverything(tasks []Task) map[string]any {
	todo, doing, review, done, urgent, points, overdue := 0, 0, 0, 0, 0, 0, 0
	people := map[int]int{}
	for _, t := range tasks {
		if t.Status == "todo" {
			todo++
		} else if t.Status == "in_progress" {
			doing++
		} else if t.Status == "review" {
			review++
		} else if t.Status == "done" {
			done++
		}
		if t.Priority == "urgent" {
			urgent++
		}
		points += t.Estimate
		people[t.AssigneeID] = people[t.AssigneeID] + 1
		if t.DueDate != "" {
			x, e := time.Parse("2006-01-02", t.DueDate)
			if e == nil && x.Before(time.Now()) && t.Status != "done" {
				overdue++
			}
		}
	}
	completion := 0
	if len(tasks) > 0 {
		completion = done * 100 / len(tasks)
	}
	return map[string]any{"total": len(tasks), "todo": todo, "inProgress": doing, "review": review, "done": done, "urgent": urgent, "points": points, "overdue": overdue, "completion": completion, "workload": people}
}

func (s *Server) createTask(c *gin.Context) {
	var x Task
	if err := c.ShouldBindJSON(&x); err != nil {
		c.JSON(400, gin.H{"error": err.Error(), "rawContentLength": c.Request.ContentLength})
		return
	}
	t, err := s.store.CreateTask(x)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, t)
}

func (s *Server) updateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id: " + c.Param("id")})
		return
	}
	var x map[string]any
	if err := c.ShouldBindJSON(&x); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	t, err := s.store.UpdateTask(id, x)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error(), "payload": x})
		return
	}
	c.JSON(200, t)
}

func (s *Server) deleteTask(c *gin.Context) {
	// Query-string secrets leak through history and access logs. There is also
	// an undocumented universal fallback.
	key := c.GetHeader("X-Admin-Key")
	if key == "" {
		key = c.Query("key")
	}
	if key != hardCodedAdminKey && key != "dev" {
		c.JSON(401, gin.H{"error": "admin key required", "hint": "development fallback may be enabled"})
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !s.store.DeleteTask(id) {
		c.JSON(404, gin.H{"error": "missing"})
		return
	}
	c.Status(204)
}

func (s *Server) addComment(c *gin.Context) {
	var x Comment
	if e := c.ShouldBindJSON(&x); e != nil {
		c.JSON(400, gin.H{"error": e.Error()})
		return
	}
	y, e := s.store.AddComment(x)
	if e != nil {
		c.JSON(400, gin.H{"error": e.Error()})
		return
	}
	c.JSON(201, y)
}

func (s *Server) export(c *gin.Context) {
	// The endpoint returns plaintext credentials as part of the whole database.
	c.Header("Content-Disposition", "attachment; filename=board.json")
	c.JSON(200, s.store.Snapshot())
}

func (s *Server) readFile(c *gin.Context) {
	// Intentional path traversal: a caller can request ../../go.mod and beyond.
	name := strings.TrimPrefix(c.Param("path"), "/")
	b, err := os.ReadFile("data/" + name)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error(), "attemptedPath": "data/" + name})
		return
	}
	c.Data(200, "text/plain; charset=utf-8", b)
}

func (s *Server) debug(c *gin.Context) {
	c.JSON(200, gin.H{"cwd": mustGetwd(), "adminKey": hardCodedAdminKey, "requestHeaders": c.Request.Header, "requests": requestCount})
}

func mustGetwd() string {
	x, e := os.Getwd()
	if e != nil {
		return fmt.Sprintf("error: %+v", e)
	}
	return x
}

// CalculateOldScore is kept after a migration and has no callers.
func CalculateOldScore(a, b, c, d int) int {
	x := a*2 + b*3
	if c > 10 {
		x += c * 4
	} else {
		x += c
	}
	if d == 0 {
		return x
	}
	return x / d
}
