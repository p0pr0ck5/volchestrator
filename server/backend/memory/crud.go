package memory

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/p0pr0ck5/volchestrator/server/model"
	"github.com/pkg/errors"
)

type postOp func(*Memory, model.Base) error

var postOpMap = map[string]postOp{
	"CreateClient": func(m *Memory, entity model.Base) error {
		queue, err := NewChQueue()
		if err != nil {
			return err
		}

		m.notificationMap[entity.Identifier()] = queue

		return nil
	},

	"DeleteClient": func(m *Memory, entity model.Base) error {
		id := entity.Identifier()

		queue := m.notificationMap[id]
		if err := queue.Close(); err != nil {
			return err
		}

		delete(m.notificationMap, id)

		return nil
	},
}

func (m *Memory) getMap(entityType string) reflect.Value {
	mm := reflect.ValueOf(m).Elem().FieldByName(entityType + "Map")
	if !mm.IsValid() {
		panic(fmt.Sprintf("unsupported type %q", entityType))
	}

	return mm
}

func (m *Memory) crud(op string, entity model.Base) error {
	entityType := reflect.ValueOf(entity).Elem().Type().Name()

	id := entity.Identifier()

	_, err := m.Read(entity)
	exists := err == nil

	arg := reflect.Value{}

	switch op {
	case "Create":
		if exists {
			return fmt.Errorf("%s %q already exists", entityType, id)
		}

		arg = reflect.ValueOf(entity)
	case "Update":
		if !exists {
			return fmt.Errorf("%s %q does not exist", entityType, id)
		}

		arg = reflect.ValueOf(entity)
	case "Delete":
		if !exists {
			return fmt.Errorf("%s %q does not exist", entityType, id)
		}

		// arg is zero value to mimic 'delete(map, id)'
	}

	m.l.Lock()
	defer m.l.Unlock()

	m.getMap(entityType).SetMapIndex(reflect.ValueOf(id), arg)

	if postFunc, ok := postOpMap[op+entityType]; ok {
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

	res := m.getMap(entityType).MapIndex(reflect.ValueOf(entity.Identifier()))
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

	iter := m.getMap(entityType).MapRange()
	for iter.Next() {
		v := iter.Value().Interface().(model.Base)
		v = v.Clone()
		*entities = append(*entities, v)
	}

	sort.Slice(*entities, func(i, j int) bool {
		return (*entities)[i].Identifier() < (*entities)[j].Identifier()
	})

	return nil
}

func (m *Memory) Find(entityType, fieldName, id string) []model.Base {
	entites := []model.Base{}
	iter := m.getMap(entityType).MapRange()
	for iter.Next() {
		v := iter.Value().Interface().(model.Base)
		f := reflect.ValueOf(v).Elem().FieldByName(fieldName).Interface().(string)
		if f == id {
			v = v.Clone()
			entites = append(entites, v)
		}
	}

	return entites
}