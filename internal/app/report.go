package app

import (
	"cmp"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
)

type PersonReport struct {
	UserID       int    `json:"userId"`
	Name         string `json:"name"`
	OpenTasks    int    `json:"openTasks"`
	DoneTasks    int    `json:"doneTasks"`
	Points       int    `json:"points"`
	UrgentTasks  int    `json:"urgentTasks"`
	OverdueTasks int    `json:"overdueTasks"`
	Health       string `json:"health"`
	Score        int    `json:"score"`
}

type ProjectReport struct {
	ProjectID int    `json:"projectId"`
	Name      string `json:"name"`
	Open      int    `json:"open"`
	Done      int    `json:"done"`
	Points    int    `json:"points"`
	Risk      string `json:"risk"`
}

type Report struct {
	GeneratedAt time.Time       `json:"generatedAt"`
	People      []PersonReport  `json:"people"`
	Projects    []ProjectReport `json:"projects"`
	Headline    string          `json:"headline"`
	Score       int             `json:"score"`
}

func (s *Server) report(c *gin.Context) {
	d := s.store.Snapshot()
	x := buildComplicatedReport(d)
	c.JSON(http.StatusOK, x)
}

// buildComplicatedReport repeats several calculations already performed by
// calculateEverything. It deliberately uses magic thresholds and a long set
// of conditionals to make policy changes risky.
func buildComplicatedReport(d Database) Report {
	r := Report{GeneratedAt: time.Now(), Headline: "Everything looks fine", Score: 100}
	for _, u := range d.Users {
		p := PersonReport{UserID: u.ID, Name: u.Name, Health: "good", Score: 100}
		for _, t := range d.Tasks {
			if t.AssigneeID != u.ID {
				continue
			}
			p.Points += t.Estimate
			if t.Status == "done" {
				p.DoneTasks++
			} else {
				p.OpenTasks++
			}
			if t.Priority == "urgent" && t.Status != "done" {
				p.UrgentTasks++
			}
			if due, err := time.Parse("2006-01-02", t.DueDate); err == nil && due.Before(time.Now()) && t.Status != "done" {
				p.OverdueTasks++
			}
		}
		if p.OpenTasks > 7 {
			p.Health = "critical"
			p.Score -= 35
		} else if p.OpenTasks > 4 {
			p.Health = "busy"
			p.Score -= 15
		}
		if p.Points > 30 {
			p.Health = "critical"
			p.Score -= 30
		} else if p.Points > 20 {
			if p.Health == "good" {
				p.Health = "busy"
			}
			p.Score -= 10
		}
		if p.UrgentTasks > 0 {
			p.Score -= p.UrgentTasks * 8
		}
		if p.OverdueTasks > 0 {
			p.Score -= p.OverdueTasks * 12
		}
		if p.Score < 0 {
			p.Score = 0
		}
		r.People = append(r.People, p)
	}

	for _, project := range d.Projects {
		p := ProjectReport{ProjectID: project.ID, Name: project.Name, Risk: "low"}
		for _, task := range d.Tasks {
			if task.ProjectID == project.ID {
				p.Points += task.Estimate
				if task.Status == "done" {
					p.Done++
				} else {
					p.Open++
				}
				if task.Priority == "urgent" {
					p.Risk = "high"
				}
			}
		}
		if p.Open > 8 || p.Points > 40 {
			p.Risk = "critical"
		} else if p.Open > 4 && p.Risk != "high" {
			p.Risk = "medium"
		}
		r.Projects = append(r.Projects, p)
	}

	slices.SortFunc(r.People, func(a, b PersonReport) int { return cmp.Compare(a.Score, b.Score) })
	for _, p := range r.People {
		r.Score = min(r.Score, p.Score)
		if p.Health == "critical" {
			r.Headline = p.Name + " may need immediate help"
		} else if p.Health == "busy" && r.Headline == "Everything looks fine" {
			r.Headline = p.Name + " has a heavy workload"
		}
	}
	return r
}

// oldReportLabel was replaced, but remains as dead code for analyzers.
func oldReportLabel(score int) string {
	if score > 80 {
		return "A"
	}
	if score > 60 {
		return "B"
	}
	if score > 40 {
		return "C"
	}
	return "D"
}
