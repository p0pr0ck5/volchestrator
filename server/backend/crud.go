package backend

import (
	"fmt"
	"reflect"
	"time"

	"github.com/p0pr0ck5/volchestrator/fsm"
	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/model"
	"github.com/p0pr0ck5/volchestrator/server/volume"
)

func (b *Backend) Create(entity model.Base) error {
	entity.Init()

	if err := entity.Validate(); err != nil {
		return err
	}

	entityType := reflect.ValueOf(entity).Elem().Type().Name()

	now := time.Now()
	set(entity, "CreatedAt", now)
	set(entity, "UpdatedAt", now)

	switch entityType {
	case "Client":
		return b.b.CreateClient(entity.(*client.Client))
	case "Volume":
		return b.b.CreateVolume(entity.(*volume.Volume))
	default:
		return fmt.Errorf("unsupported type %q", entityType)
	}
}

func (b *Backend) Read(entity model.Base) error {
	entityType := reflect.ValueOf(entity).Elem().Type().Name()

	var f model.Base
	var err error

	switch entityType {
	case "Client":
		c := entity.(*client.Client)
		f, err = b.b.ReadClient(c.ID)
	case "Volume":
		v := entity.(*volume.Volume)
		f, err = b.b.ReadVolume(v.ID)
	default:
		return fmt.Errorf("unsupported type %q", entityType)
	}

	if err != nil {
		return err
	}

	copy(f, entity, true)

	return nil
}

func (b *Backend) Update(entity model.Base) error {
	entityType := reflect.ValueOf(entity).Elem().Type().Name()

	var f model.Base
	var err error

	switch entityType {
	case "Client":
		c := entity.(*client.Client)
		f, err = b.b.ReadClient(c.ID)
	case "Volume":
		v := entity.(*volume.Volume)
		f, err = b.b.ReadVolume(v.ID)
	default:
		return fmt.Errorf("unsupported type %q", entityType)
	}

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

	if err := entity.F().Transition(s); err != nil {
		return err
	}

	set(entity, "UpdatedAt", time.Now())

	switch entityType {
	case "Client":
		return b.b.UpdateClient(entity.(*client.Client))
	case "Volume":
		return b.b.UpdateVolume(entity.(*volume.Volume))
	default:
		return fmt.Errorf("unsupported type %q", entityType)
	}
}

func (b *Backend) Delete(entity model.Base) error {
	entityType := reflect.ValueOf(entity).Elem().Type().Name()

	var err error
	var i int

	switch entityType {
	case "Client":
		c := entity.(*client.Client)
		_, err = b.b.ReadClient(c.ID)
		i = int(client.Deleting)
	case "Volume":
		v := entity.(*volume.Volume)
		_, err = b.b.ReadVolume(v.ID)
		i = int(volume.Deleting)
	default:
		return fmt.Errorf("unsupported type %q", entityType)
	}

	if err != nil {
		return err
	}

	// set Deleting status, and update
	reflect.ValueOf(entity).Elem().FieldByName("Status").SetInt(int64(i))

	if err := b.Update(entity); err != nil {
		return err
	}

	switch entityType {
	case "Client":
		return b.b.DeleteClient(entity.(*client.Client))
	case "Volume":
		return b.b.DeleteVolume(entity.(*volume.Volume))
	default:
		return fmt.Errorf("unsupported type %q", entityType)
	}
}

func (b *Backend) List(entityType string, entities *[]model.Base) error {
	switch entityType {
	case "client":
		clients, err := b.b.ListClients()
		if err != nil {
			return err
		}

		for _, client := range clients {
			*entities = append(*entities, client)
		}

		return nil
	case "volume":
		volumes, err := b.b.ListVolumes()
		if err != nil {
			return err
		}

		for _, volume := range volumes {
			*entities = append(*entities, volume)
		}

		return nil
	default:
		return fmt.Errorf("unsupported type %q", entityType)
	}
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
