package backend

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/model"
)

type processfunc func(string, int, reflect.StructField) error

func (b *Backend) Create(entity model.Base) error {
	entity.Init()

	if err := entity.Validate(); err != nil {
		return err
	}

	err := b.processModel(
		entity,
		func(tag string, i int, field reflect.StructField) error {
			if tag != "required" {
				return nil
			}

			fieldVal := reflect.ValueOf(entity).Elem().Field(i)

			if fieldVal.IsZero() {
				return fmt.Errorf("validate error (required): %q", fieldVal)
			}

			return nil
		},
		func(tag string, i int, field reflect.StructField) error {
			if !strings.Contains(tag, "depends") {
				return nil
			}

			fieldVal := reflect.ValueOf(entity).Elem().Field(i)

			v := strings.Split(tag, "=")[1]
			e, key := strings.Split(v, ":")[0], strings.Split(v, ":")[1]

			ee := b.b.Find(e, key, fieldVal.Interface().(string))
			if len(ee) == 0 {
				return fmt.Errorf("missing reference %v:%v", entity, key)
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	now := time.Now()
	set(entity, "CreatedAt", now)
	set(entity, "UpdatedAt", now)

	return b.b.Create(entity)
}

func (b *Backend) Read(entity model.Base) error {
	ff, err := b.b.Read(entity)
	if err != nil {
		return err
	}
	f := ff.Clone()

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

	err = b.processModel(
		entity,
		func(tag string, i int, field reflect.StructField) error {
			if tag != "immutable" {
				return nil
			}

			fieldVal := reflect.ValueOf(entity).Elem().Field(i).Interface()
			newFieldVal := reflect.ValueOf(f).Elem().Field(i).Interface()

			if fieldVal != newFieldVal {
				return fmt.Errorf("validate error (immutable): %q != %q", fieldVal, newFieldVal)
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	var s fsm.State
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

	err := b.processModel(
		entity,
		func(tag string, i int, field reflect.StructField) error {
			if !strings.Contains(tag, "reference") {
				return nil
			}

			fieldVal := reflect.ValueOf(entity).Elem().Field(i)
			v := strings.Split(tag, "=")[1]
			e, key := strings.Split(v, ":")[0], strings.Split(v, ":")[1]

			ee := b.b.Find(e, key, fieldVal.Interface().(string))
			if ee == nil {
				return fmt.Errorf("missing reference %v:%v", entity, key)
			}

			for _, dependent := range ee {
				if err := b.Delete(dependent); err != nil {
					return fmt.Errorf("unable to delete referenced entity: %w", err)
				}
			}

			return nil
		},
	)
	if err != nil {
		return err
	}

	return b.b.Delete(entity)
}

func (b *Backend) List(entityType string, entities *[]model.Base) error {
	return b.b.List(entityType, entities)
}

func (b *Backend) processModel(entity model.Base, processors ...processfunc) error {
	entityType := reflect.TypeOf(entity).Elem()
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)

		modelTags := strings.Split(field.Tag.Get("model"), ",")

		for _, s := range modelTags {
			for _, f := range processors {
				if err := f(s, i, field); err != nil {
					return err
				}
			}
		}
	}

	return nil
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
