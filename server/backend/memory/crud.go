package memory

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/model"
	"github.com/pkg/errors"
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

func (m *Memory) crud(op string, entity model.Base) error {
	entityType := reflect.ValueOf(entity).Elem().Type().Name()

	id := entity.Identifier()

	_, err := m.Read(entity)
	exists := err == nil

	dMap := reflect.ValueOf(m).Elem().FieldByName(entityType + "Map")

	m.l.Lock()
	defer m.l.Unlock()

	switch op {
	case "Create":
		if exists {
			return fmt.Errorf("%s %q already exists", entityType, id)
		}

		dMap.SetMapIndex(reflect.ValueOf(id), reflect.ValueOf(entity))
	case "Update":
		if !exists {
			return fmt.Errorf("%s %q already exists", entityType, id)
		}

		dMap.SetMapIndex(reflect.ValueOf(id), reflect.ValueOf(entity))
	case "Delete":
		if !exists {
			return fmt.Errorf("%s %q already exists", entityType, id)
		}

		dMap.SetMapIndex(reflect.ValueOf(id), reflect.Value{})
	}

	if postFunc, ok := post_op_map[op+entityType]; ok {
		postFunc(m, entity)
	}

	return nil
}

func (m *Memory) Create(entity model.Base) error {
	return m.crud("Create", entity)
}

func (m *Memory) Read(entity model.Base) (model.Base, error) {
	m.l.RLock()
	defer m.l.RUnlock()

	entityType := reflect.ValueOf(entity).Elem().Type().Name()

	id := entity.Identifier()
	dMap := reflect.ValueOf(m).Elem().FieldByName(entityType + "Map")

	res := dMap.MapIndex(reflect.ValueOf(id))
	if !res.IsValid() {
		return nil, errors.New("not found")
	}

	return res.Interface().(model.Base), nil
}

func (m *Memory) Update(entity model.Base) error {
	return m.crud("Update", entity)
}

func (m *Memory) Delete(entity model.Base) error {
	return m.crud("Delete", entity)
}

func (m *Memory) List(entityType string, entities *[]model.Base) error {
	m.l.RLock()
	defer m.l.RUnlock()

	dMap := reflect.ValueOf(m).Elem().FieldByName(strings.Title(entityType) + "Map")

	iter := dMap.MapRange()
	for iter.Next() {
		v := iter.Value().Interface().(model.Base)
		*entities = append(*entities, v)
	}

	sort.Slice(*entities, func(i, j int) bool {
		return (*entities)[i].Identifier() < (*entities)[j].Identifier()
	})

	return nil
}
