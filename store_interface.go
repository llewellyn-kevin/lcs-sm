package main

import(
  "fmt"
  "reflect"

  "github.com/gomodule/redigo/redis"
)

// HashToStruct uses a store connection and a key generated from a
// hashName and hashId (using the format {hashName}:{hashId}) to
// assign all the values of a struct. 
func HashToStruct(store redis.Conn, key string, target interface{}) error {
  targetValue := reflect.ValueOf(target).Elem()
  targetType := targetValue.Type()

  id := targetValue.FieldByName("Id").Int()
  key = fmt.Sprintf("%s:%d", key, id)

  for i := 0; i < targetType.NumField(); i++ {
    field := targetType.Field(i)
    tag := field.Tag.Get("hash-key")
    if tag != "" { // This field is a key in the hash
      value := targetValue.FieldByName(field.Name)
      switch value.Kind() {
        case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
          if result, err := redis.Int64(store.Do("HGET", key, tag)); err != nil {
            return err
          } else {
            value.SetInt(result)
          }
        case reflect.Bool:
          if result, err := redis.Bool(store.Do("HGET", key, tag)); err != nil {
            return err
          } else {
            value.SetBool(result)
          }
        case reflect.Float64, reflect.Float32:
          if result, err := redis.Float64(store.Do("HGET", key, tag)); err != nil {
            return err
          } else {
            value.SetFloat(result)
          }
        case reflect.String:
          if result, err := redis.String(store.Do("HGET", key, tag)); err != nil {
            return err
          } else {
            value.SetString(result)
          }
        default:
          return fmt.Errorf("Could not identify type of field: %s", field.Name)
      }
    }
  }

  return nil
}
