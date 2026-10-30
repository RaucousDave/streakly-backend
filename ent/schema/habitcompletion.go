package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// HabitCompletion holds the schema definition for the HabitCompletion entity.
type HabitCompletion struct {
	ent.Schema
}

func (HabitCompletion) Mixin() []ent.Mixin { return []ent.Mixin{mixin.Time{}} }

// Fields of the HabitCompletion.
func (HabitCompletion) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("user_id", uuid.UUID{}),
		field.UUID("habit_id", uuid.UUID{}),
		field.String("completed_date"),
		field.Time("completed_at_utc"),
		field.Int("qty_achieved"),
		field.Bool("is_mvt_met").Default(false),
		field.Bool("is_target_met").Default(false),
		field.String("idempotency_key").Unique(),
	}
}

// Edges of the HabitCompletion.
func (HabitCompletion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("habit_completions").Field("user_id").Unique().Required(),
		edge.From("habit", Habit.Type).Ref("habit_completions").Field("habit_id").Unique().Required(),
	}
}

func (HabitCompletion) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("habit_id", "completed_date").Unique(),
	}
}
