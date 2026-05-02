package dto

type DashboardStatsQuery struct {
	ProjectID uint `form:"project_id" binding:"required,gt=0"`
}

type TasksPerUserResponse struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	Count    int64  `json:"count"`
}

type DashboardStatsResponse struct {
	ProjectID     uint                   `json:"project_id"`
	TotalTasks    int64                  `json:"total_tasks"`
	OverdueTasks  int64                  `json:"overdue_tasks"`
	TasksByStatus map[string]int64       `json:"tasks_by_status"`
	TasksPerUser  []TasksPerUserResponse `json:"tasks_per_user"`
}
