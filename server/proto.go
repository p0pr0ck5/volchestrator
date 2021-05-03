package server

import (
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/p0pr0ck5/volchestrator/server/client"
	"github.com/p0pr0ck5/volchestrator/server/notification"
	"github.com/p0pr0ck5/volchestrator/server/volume"
	"github.com/p0pr0ck5/volchestrator/svc"
)

func toProto(from interface{}) interface{} {
	var to interface{}
	var val reflect.Value
	var ctx string

	switch from.(type) {
	case *client.Client:
		to = &svc.Client{}
		val = reflect.ValueOf(from.(*client.Client)).Elem()
		ctx = "Client"
	case *volume.Volume:
		to = &svc.Volume{}
		val = reflect.ValueOf(from.(*volume.Volume)).Elem()
		ctx = "Volume"
	case *notification.Notification:
		to = &svc.Notification{}
		val = reflect.ValueOf(from.(*notification.Notification)).Elem()
		ctx = "Notification"
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		var fieldName string

		// if this is the ID, its a special case
		switch field.Name {
		case "ID":
			fieldName = ctx + "Id"
		default:
			fieldName = field.Name
		}

		toField := reflect.ValueOf(to).Elem().FieldByName(fieldName)
		if toField.CanSet() {
			switch field.Type.Kind() {
			case reflect.String:
				toField.SetString(val.Field(i).String())
			case reflect.Int:
				toField.SetInt(val.Field(i).Int())
			default:
				switch field.Type.String() {
				case "time.Time":
					seconds := val.Field(i).MethodByName("Unix").Call([]reflect.Value{})[0].Int()
					nano := val.Field(i).MethodByName("Nanosecond").Call([]reflect.Value{})[0].Int()

					// zero time, just leave it unset
					if seconds != -62135596800 {
						tp, _ := ptypes.TimestampProto(time.Unix(seconds, nano))
						toField.Set(reflect.ValueOf(tp))
					}
				default:
					panic("unsupported type " + field.Type.String())
				}
			}
		}
	}

	return to
}

func toStruct(from interface{}) interface{} {
	var to interface{}
	var val reflect.Value
	var ctx string

	switch from.(type) {
	case *svc.Volume:
		to = &volume.Volume{}
		val = reflect.ValueOf(from.(*svc.Volume)).Elem()
		ctx = "Volume"
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		var fieldName string

		switch field.Name {
		case ctx + "Id":
			fieldName = "ID"
		default:
			fieldName = field.Name
		}

		toField := reflect.ValueOf(to).Elem().FieldByName(fieldName)
		if toField.CanSet() {
			switch field.Type.Kind() {
			case reflect.String:
				toField.SetString(val.Field(i).String())
			case reflect.Int, reflect.Int32:
				toField.SetInt(val.Field(i).Int())
			default:
				switch field.Type.String() {
				case "*timestamppb.Timestamp":
					// TODO
				default:
					panic("unsupported type " + field.Type.String())
				}
			}
		}
	}

	return to
}
