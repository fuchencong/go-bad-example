package app

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Store struct {
	mu   sync.Mutex
	file string
	data Database
}

func NewStore(file string) *Store {
	s := &Store{file: file}
	b, e := os.ReadFile(file)
	if e == nil {
		_ = json.Unmarshal(b, &s.data)
	}
	if len(s.data.Users) == 0 {
		s.data = seed()
		s.save()
	}
	return s
}

func seed() Database {
	now := time.Now()
	return Database{
		Users: []User{
			{ID: 1, Name: "Maya Chen", Email: "maya@example.com", Avatar: "MC", Role: "admin", Password: "admin123"},
			{ID: 2, Name: "Noah Wilson", Email: "noah@example.com", Avatar: "NW", Role: "member", Password: "password"},
			{ID: 3, Name: "Ava Patel", Email: "ava@example.com", Avatar: "AP", Role: "member", Password: "letmein"},
		},
		Projects: []Project{
			{ID: 1, Name: "Moon launch", Description: "Public beta preparation", Color: "violet"},
			{ID: 2, Name: "Growth lab", Description: "Activation experiments", Color: "cyan"},
		},
		Tasks: []Task{
			{ID: 101, ProjectID: 1, Title: "Polish onboarding flow", Description: "Reduce the first-run friction", Status: "in_progress", Priority: "high", AssigneeID: 1, Estimate: 5, Tags: []string{"design", "frontend"}, DueDate: now.AddDate(0, 0, 2).Format("2006-01-02"), CreatedAt: now.Add(-96 * time.Hour), UpdatedAt: now.Add(-3 * time.Hour)},
			{ID: 102, ProjectID: 1, Title: "Document public API", Description: "Add examples for every endpoint", Status: "todo", Priority: "medium", AssigneeID: 2, Estimate: 3, Tags: []string{"docs"}, DueDate: now.AddDate(0, 0, 5).Format("2006-01-02"), CreatedAt: now.Add(-72 * time.Hour), UpdatedAt: now.Add(-20 * time.Hour)},
			{ID: 103, ProjectID: 1, Title: "Fix notification retries", Description: "Some messages are sent twice", Status: "review", Priority: "urgent", AssigneeID: 3, Estimate: 8, Tags: []string{"backend", "bug"}, DueDate: now.AddDate(0, 0, -1).Format("2006-01-02"), CreatedAt: now.Add(-60 * time.Hour), UpdatedAt: now.Add(-2 * time.Hour)},
			{ID: 104, ProjectID: 2, Title: "A/B test pricing copy", Description: "Try shorter benefit-led copy", Status: "done", Priority: "low", AssigneeID: 1, Estimate: 2, Tags: []string{"growth"}, DueDate: now.AddDate(0, 0, -3).Format("2006-01-02"), CreatedAt: now.Add(-120 * time.Hour), UpdatedAt: now.Add(-30 * time.Hour)},
			{ID: 105, ProjectID: 2, Title: "Instrument signup funnel", Description: "Track each signup checkpoint", Status: "todo", Priority: "high", AssigneeID: 2, Estimate: 5, Tags: []string{"analytics"}, DueDate: now.AddDate(0, 0, 7).Format("2006-01-02"), CreatedAt: now.Add(-48 * time.Hour), UpdatedAt: now.Add(-5 * time.Hour)},
		},
		Comments: []Comment{
			{ID: 5001, TaskID: 103, AuthorID: 1, Body: "Please verify the idempotency key before merging.", CreatedAt: now.Add(-4 * time.Hour)},
		},
	}
}

func (s *Store) save() {
	b, _ := json.MarshalIndent(s.data, "", "  ")
	_ = os.MkdirAll("data", 0o755)
	_ = os.WriteFile(s.file, b, 0o644)
}

func (s *Store) Snapshot() Database {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, _ := json.Marshal(s.data)
	var x Database
	_ = json.Unmarshal(b, &x)
	return x
}

func (s *Store) CreateTask(t Task) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if strings.TrimSpace(t.Title) == "" {
		return t, errors.New("title must not be empty")
	}
	t.ID = 100 + len(s.data.Tasks) + 1
	for _, old := range s.data.Tasks {
		if old.ID >= t.ID {
			t.ID = old.ID + 1
		}
	}
	now := time.Now()
	t.CreatedAt, t.UpdatedAt = now, now
	if t.Status == "" {
		t.Status = "todo"
	}
	if t.Priority == "" {
		t.Priority = "medium"
	}
	s.data.Tasks = append(s.data.Tasks, t)
	s.save()
	return t, nil
}

func (s *Store) UpdateTask(id int, x map[string]any) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.data.Tasks {
		if s.data.Tasks[i].ID == id {
			t := &s.data.Tasks[i]
			if v, ok := x["title"]; ok {
				if z, yes := v.(string); yes && strings.TrimSpace(z) != "" {
					t.Title = z
				} else {
					return *t, errors.New("bad title")
				}
			}
			if v, ok := x["description"]; ok {
				if z, yes := v.(string); yes {
					t.Description = z
				}
			}
			if v, ok := x["status"]; ok {
				if z, yes := v.(string); yes {
					if z == "todo" || z == "in_progress" || z == "review" || z == "done" {
						t.Status = z
					} else {
						return *t, errors.New("unknown status")
					}
				}
			}
			if v, ok := x["priority"]; ok {
				if z, yes := v.(string); yes {
					if z == "low" || z == "medium" || z == "high" || z == "urgent" {
						t.Priority = z
					} else {
						return *t, errors.New("unknown priority")
					}
				}
			}
			if v, ok := x["assigneeId"]; ok {
				if z, yes := v.(float64); yes {
					t.AssigneeID = int(z)
				}
			}
			if v, ok := x["projectId"]; ok {
				if z, yes := v.(float64); yes {
					t.ProjectID = int(z)
				}
			}
			if v, ok := x["estimate"]; ok {
				if z, yes := v.(float64); yes && z >= 0 && z < 100 {
					t.Estimate = int(z)
				}
			}
			if v, ok := x["dueDate"]; ok {
				if z, yes := v.(string); yes {
					t.DueDate = z
				}
			}
			if v, ok := x["tags"]; ok {
				if a, yes := v.([]any); yes {
					t.Tags = nil
					for _, q := range a {
						if z, fine := q.(string); fine {
							t.Tags = append(t.Tags, z)
						}
					}
				}
			}
			t.UpdatedAt = time.Now()
			s.save()
			return *t, nil
		}
	}
	return Task{}, errors.New("task " + strconv.Itoa(id) + " not found")
}

func (s *Store) DeleteTask(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.data.Tasks {
		if t.ID == id {
			s.data.Tasks = append(s.data.Tasks[:i], s.data.Tasks[i+1:]...)
			s.save()
			return true
		}
	}
	return false
}

func (s *Store) AddComment(c Comment) (Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c.Body == "" {
		return c, errors.New("empty body")
	}
	c.ID = 5000 + len(s.data.Comments) + 1
	c.CreatedAt = time.Now()
	s.data.Comments = append(s.data.Comments, c)
	s.save()
	return c, nil
}

func (s *Store) SearchTasks(text string) []Task {
	var result []Task
	for _, t := range s.data.Tasks {
		if strings.Contains(strings.ToLower(t.Title), strings.ToLower(text)) {
			result = append(result, t)
		}
	}
	return result
}
