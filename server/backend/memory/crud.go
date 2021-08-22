package memory

import (
	"fmt"
	"reflect"

	"github.com/p0pr0ck5/volchestrator/server/model"
)

func (m *Memory) crud(op string, entity model.Base) error {
	entityType := reflect.ValueOf(entity).Elem().Type().Name()

	fnName := op + entityType

	ff := reflect.ValueOf(m).MethodByName(fnName)

	if !ff.IsValid() {
		return fmt.Errorf("unsupported type %q", entityType)
	}

	res := ff.Call([]reflect.Value{
		reflect.ValueOf(entity),
	})

	if err, ok := res[0].Interface().(error); ok {
		return err
	} else {
		return nil
	}
}

func (m *Memory) Create(entity model.Base) error {
	return m.crud("Create", entity)
}

func (m *Memory) Read(entity model.Base) (model.Base, error) {
	entityType := reflect.ValueOf(entity).Elem().Type().Name()

	fnName := "Read" + entityType

	ff := reflect.ValueOf(m).MethodByName(fnName)

	if !ff.IsValid() {
		return nil, fmt.Errorf("unsupported type %q", entityType)
	}

	id := reflect.ValueOf(entity).Elem().FieldByName("ID")

	res := ff.Call([]reflect.Value{
		id,
	})

	if err, ok := res[1].Interface().(error); ok {
		return nil, err
	} else {
		return res[0].Interface().(model.Base), nil
	}
}

func (m *Memory) Update(entity model.Base) error {
	return m.crud("Update", entity)
}

func (m *Memory) Delete(entity model.Base) error {
	return m.crud("Delete", entity)
}

func (m *Memory) List(entityType string, entities *[]model.Base) error {
	switch entityType {
	case "client":
		clients, err := m.ListClients()
		if err != nil {
			return err
		}

		for _, client := range clients {
			*entities = append(*entities, client)
		}

		return nil
	case "volume":
		volumes, err := m.ListVolumes()
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
