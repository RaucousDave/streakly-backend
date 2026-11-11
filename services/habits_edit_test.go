package services

import (
	"errors"
	"testing"

	"github.com/RaucousDave/streakly/dto"
	"github.com/RaucousDave/streakly/ent"
	"github.com/RaucousDave/streakly/ent/habit"
)

func TestPrepareHabitEdit(t *testing.T) {
	unit, title := "pages", "  Read daily  "
	target, minimum, tooSmall, zero := 20, 5, 3, 0
	boolean, quantitative, invalid := dto.HabitTypeBoolean, dto.HabitTypeQuantitative, dto.HabitType("INVALID")
	blank := "  "
	current := ent.Habit{Title: "Read", Type: habit.TypeQUANTITATIVE, Unit: &unit, TargetQty: &target, MvtQty: &minimum}
	tests := []struct {
		name      string
		current   ent.Habit
		input     dto.EditHabitDto
		wantError bool
	}{
		{"rename preserves quantities", current, dto.EditHabitDto{Title: &title}, false},
		{"empty patch", current, dto.EditHabitDto{}, true},
		{"blank title", current, dto.EditHabitDto{Title: &blank}, true},
		{"blank unit", current, dto.EditHabitDto{Unit: &blank}, true},
		{"invalid type", current, dto.EditHabitDto{Type: &invalid}, true},
		{"target below saved minimum", current, dto.EditHabitDto{TargetQty: &tooSmall}, true},
		{"zero minimum", current, dto.EditHabitDto{MVTQty: &zero}, true},
		{"switch to boolean", current, dto.EditHabitDto{Type: &boolean}, false},
		{"incomplete quantitative conversion", ent.Habit{Title: "Read", Type: habit.TypeBOOLEAN}, dto.EditHabitDto{Type: &quantitative}, true},
		{"complete quantitative conversion", ent.Habit{Title: "Read", Type: habit.TypeBOOLEAN}, dto.EditHabitDto{Type: &quantitative, TargetQty: &target, MVTQty: &minimum, Unit: &unit}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next, err := prepareHabitEdit(tt.current, tt.input)
			if tt.wantError {
				if !errors.Is(err, ErrInvalidHabitEdit) {
					t.Fatalf("expected validation error, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if next.Type == habit.TypeBOOLEAN && (next.Unit != nil || next.TargetQty != nil || next.MvtQty != nil) {
				t.Fatal("boolean habit retained quantitative fields")
			}
			if tt.input.Title != nil && next.Title != "Read daily" {
				t.Fatalf("title was not trimmed: %q", next.Title)
			}
			if next.Type == habit.TypeQUANTITATIVE && (*next.TargetQty != target || *next.MvtQty != minimum || *next.Unit != unit) {
				t.Fatal("quantitative fields were not preserved")
			}
		})
	}
}
