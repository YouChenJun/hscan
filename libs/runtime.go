package libs

import "time"

// runtime metadata
type RuntimeInfo struct {
	Target    string
	Workspace string
	CreatedAt time.Time `gorm:"column:created_at;type:datetime(0);autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime(0);autoUpdateTime" json:"updated_at,omitempty"`
}
