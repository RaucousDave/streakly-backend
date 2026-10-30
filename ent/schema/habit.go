package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Habit holds the schema definition for the Habit entity.
type Habit struct {
	ent.Schema
}

func (Habit) Mixin() []ent.Mixin { return []ent.Mixin{mixin.Time{}} }

// Fields of the Habit.
func (Habit) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("user_id", uuid.UUID{}),
		field.String("title"),
		field.Enum("type").Values("BOOLEAN", "QUANTITATIVE"),
		field.Int("target_qty").Optional().Nillable(),
		field.Int("mvt_qty").Optional().Nillable(),
		field.String("unit").Optional().Nillable(),
		field.Bool("is_archived").Default(false),
		field.Time("archived_at").Optional().Nillable(),
	}
}

// Edges of the Habit.
func (Habit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("habits").Field("user_id").Unique().Required(),
		edge.To("freeze_logs", FreezeLog.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("habit_completions", HabitCompletion.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
