package memory

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/model"
)

type post_op func(*Memory, model.Base) error

var post_op_map = map[string]post_op{
	"CreateClient": func(m *Memory, entity model.Base) error {
		queue, err := NewChQueue()
		if err != nil {
			return err
		}

		m.notificationMap[entity.(*client.Client).ID] = queue

		return nil
	},

	"DeleteClient": func(m *Memory, entity model.Base) error {
		client := entity.(*client.Client)

		queue := m.notificationMap[client.ID]
		if err := queue.Close(); err != nil {
			return err
		}
		delete(m.notificationMap, client.ID)

		return nil
	},
}

func (m *Memory) read(id, entityType string) (interface{}, error) {
	entity, exists := m.getMap(entityType).Get(id)
	if !exists {
		return nil, fmt.Errorf("%s %q not found", entityType, id)
	}

	return entity, nil
}

func (m *Memory) crud(op string, entity model.Base) error {
	entityType := reflect.ValueOf(entity).Elem().Type().Name()

	id := reflect.ValueOf(entity).Elem().FieldByName("ID").Interface().(string)

	_, err := m.read(id, strings.ToLower(entityType))
	exists := err == nil

	dataMap := reflect.ValueOf(m).Elem().FieldByName(entityType + "Map")
	var fn string
	var args []reflect.Value

	switch op {
	case "Create":
		if exists {
			return fmt.Errorf("%s %q already exists", entityType, id)
		}

		fn = "Set"
		args = []reflect.Value{reflect.ValueOf(id), reflect.ValueOf(entity)}
	case "Update":
		if !exists {
			return fmt.Errorf("%s %q already exists", entityType, id)
		}

		fn = "Set"
		args = []reflect.Value{reflect.ValueOf(id), reflect.ValueOf(entity)}
	case "Delete":
		if !exists {
			return fmt.Errorf("%s %q already exists", entityType, id)
		}

		fn = "Delete"
		args = []reflect.Value{reflect.ValueOf(id)}
	}

	ff := dataMap.MethodByName(fn)
	if !ff.IsValid() {
		return fmt.Errorf("unsupported type %q", entityType)
	}

	ff.Call(args)

	if postFunc, ok := post_op_map[op+entityType]; ok {
		postFunc(m, entity)
	}

	return nil
}

func (m *Memory) Create(entity model.Base) error {
	return m.crud("Create", entity)
}

func (m *Memory) Read(entity model.Base) (model.Base, error) {
	entityType := reflect.ValueOf(entity).Elem().Type().Name()

	id := reflect.ValueOf(entity).Elem().FieldByName("ID").Interface().(string)
	res, err := m.read(id, strings.ToLower(entityType))

	if err != nil {
		return nil, err
	} else {
		return res.(model.Base), err
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
