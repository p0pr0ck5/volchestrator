package backend

import (
	"reflect"
	"time"

	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/model"
)

func (b *Backend) Create(entity model.Base) error {
	entity.Init()

	if err := entity.Validate(); err != nil {
		return err
	}

	now := time.Now()
	set(entity, "CreatedAt", now)
	set(entity, "UpdatedAt", now)

	return b.b.Create(entity)
}

func (b *Backend) Read(entity model.Base) error {
	f, err := b.b.Read(entity)
	if err != nil {
		return err
	}

	copy(f, entity, true)

	return nil
}

func (b *Backend) Update(entity model.Base) error {
	f, err := b.b.Read(entity)
	if err != nil {
		return err
	}

	copy(f, entity, false)

	if err := entity.Validate(); err != nil {
		return err
	}

	if err := f.ValidateTransition(entity); err != nil {
		return err
	}

	var s fsm.State = 0
	i := reflect.ValueOf(entity).Elem().FieldByName("Status")
	reflect.ValueOf(&s).Elem().Set(i)

	if err := entity.F().Transition(s, entity); err != nil {
		return err
	}

	set(entity, "UpdatedAt", time.Now())

	return b.b.Update(entity)
}

func (b *Backend) Delete(entity model.Base) error {
	entity.SetStatus("Deleting")

	if err := b.Update(entity); err != nil {
		return err
	}

	return b.b.Delete(entity)
}

func (b *Backend) List(entityType string, entities *[]model.Base) error {
	return b.b.List(entityType, entities)
}

func copy(src, dest interface{}, fullMerge bool) {
	v := reflect.ValueOf(src).Elem()
	for i := 0; i < v.NumField(); i++ {
		sourceField := v.Field(i)
		destField := reflect.ValueOf(dest).Elem().Field(i)

		if destField.CanSet() {
			if destField.IsZero() || fullMerge {
				destField.Set(sourceField)
			}
		}
	}
}

func set(entity interface{}, field string, value interface{}) {
	f := reflect.ValueOf(entity).Elem().FieldByName(field)

	if f.CanSet() {
		f.Set(reflect.ValueOf(value))
	}
}
