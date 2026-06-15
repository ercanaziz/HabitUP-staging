package habit

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Habit struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"userId" json:"userId"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
}

type CheckRecord struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	HabitID   primitive.ObjectID `bson:"habitId" json:"habitId"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Date      string             `bson:"date" json:"date"` // "2024-06-15" formatı
	CheckedAt time.Time          `bson:"checkedAt" json:"checkedAt"`
}

type HabitStats struct {
	CurrentStreak  int     `json:"currentStreak"`
	LongestStreak  int     `json:"longestStreak"`
	CompletionRate float64 `json:"completionRate"`
	TotalChecks    int     `json:"totalChecks"`
}

type CreateHabitRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateHabitRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
