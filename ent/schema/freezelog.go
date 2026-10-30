package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// FreezeLog holds the schema definition for the FreezeLog entity.
type FreezeLog struct {
	ent.Schema
}

func (FreezeLog) Mixin() []ent.Mixin { return []ent.Mixin{mixin.Time{}} }

// Fields of the FreezeLog.
func (FreezeLog) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Enum("reason").Values("WEEKLY_GRANT", "PERFORMANCE_REWARD", "MANUAL_SPEND"),
		field.UUID("habit_id", uuid.UUID{}).Optional().Nillable(),
		field.Int("amount"),
		field.UUID("user_id", uuid.UUID{}),
		field.String("target_date"),
	}
}

// Edges of the FreezeLog.
func (FreezeLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("freeze_logs").Field("user_id").Unique().Required(),
		edge.From("habit", Habit.Type).Ref("freeze_logs").Field("habit_id").Unique(),
	}
}
