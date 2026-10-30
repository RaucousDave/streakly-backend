package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin { return []ent.Mixin{mixin.Time{}} }

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("email").Unique().NotEmpty(),
		field.String("google_id").Optional().Nillable().Unique().Sensitive(),
		field.String("password_hash").Optional().Nillable().Sensitive(),
		field.Enum("auth_provider").Values("GOOGLE", "EMAIL"),
		field.String("display_name").NotEmpty(),
		field.String("timezone"),
		field.String("avatar_url").Optional().Nillable(),
		field.Int("freeze_balance").Default(1),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("habits", Habit.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("freeze_logs", FreezeLog.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("habit_completions", HabitCompletion.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
