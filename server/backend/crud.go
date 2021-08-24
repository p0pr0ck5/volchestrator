package backend

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/model"
)

func (b *Backend) Create(entity model.Base) error {
	entity.Init()

	if err := entity.Validate(); err != nil {
		return err
	}

	// struct validation
	entityType := reflect.TypeOf(entity).Elem()
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		fieldVal := reflect.ValueOf(entity).Elem().Field(i)

		modelTags := strings.Split(field.Tag.Get("model"), ",")

		// required check
		for _, s := range modelTags {
			if s == "required" {
				if fieldVal.IsZero() {
					return fmt.Errorf("validate error (required): %q", fieldVal)
				}
			}

			if strings.Contains(s, "reference") {
				v := strings.Split(s, "=")[1]
				e, key := strings.Split(v, ":")[0], strings.Split(v, ":")[1]

				ee := b.b.Find(e, fieldVal.Interface().(string))
				if ee == nil {
					return fmt.Errorf("missing reference %v:%v", entity, key)
				}
			}
		}
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

	// struct validation
	entityType := reflect.TypeOf(entity).Elem()
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		fieldVal := reflect.ValueOf(entity).Elem().Field(i).Interface()
		newFieldVal := reflect.ValueOf(f).Elem().Field(i).Interface()

		modelTags := strings.Split(field.Tag.Get("model"), ",")

		// immutable check
		for _, s := range modelTags {
			if s != "immutable" {
				continue
			}

			if fieldVal != newFieldVal {
				return fmt.Errorf("validate error (immutable): %q != %q", fieldVal, newFieldVal)
			}
		}
	}

	var s fsm.State
	i := reflect.ValueOf(entity).Elem().FieldByName("Status")
	reflect.ValueOf(&s).Elem().Set(i)

	if can := f.F().Can(s); !can {
		return fmt.Errorf("invalid state transition %q", s)
	}

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
