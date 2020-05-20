package main

import(
  "fmt"
  "reflect"

  "github.com/gertd/go-pluralize"
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
    tag := field.Tag.Get("redii")
    if tag != "" && tag != "pk" { // This field is a key in the hash
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

// CreateHash creates a hash resource in a redis store based on the
// data in a struct.
func CreateHash(store redis.Conn, key string, target interface{}) error {
  // Generate keys
  pluralize := pluralize.NewClient()
  pkey := pluralize.Plural(key) // plural key for list of ids
  ikey := fmt.Sprintf("next-%s-id", key) // incrementer key for next id

  // Reflect on parts of struct
  targetValue := reflect.ValueOf(target).Elem()
  targetType := targetValue.Type()

  // Increment ID for resource
  id, err := redis.Int(store.Do("INCR", ikey))
  if err != nil {
    return err
  }

  // Get key for hash
  key = fmt.Sprintf("%s:%d", key, id)

  // Create hash at key
  for i := 0; i < targetType.NumField(); i++ {
    field := targetType.Field(i)
    tag := field.Tag.Get("redii")

    if tag == "pk" {
      value := targetValue.FieldByName(field.Name)
      value.SetInt(int64(id))
    } else {
      if tag != "" { // This field is a key in the hash
        value := targetValue.FieldByName(field.Name)
        switch value.Kind() {
          case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
            unpacked := value.Uint()
            if _, err := store.Do("HSET", key, tag, unpacked); err != nil {
              return err
            }
          case reflect.Bool:
            unpacked := value.Bool()
            if _, err := store.Do("HSET", key, tag, unpacked); err != nil {
              return err
            }
          case reflect.Float64, reflect.Float32:
            unpacked := value.Float()
            if _, err := store.Do("HSET", key, tag, unpacked); err != nil {
              return err
            }
          case reflect.String:
            unpacked := value.String()
            if _, err := store.Do("HSET", key, tag, unpacked); err != nil {
              return err
            }
          default:
            return fmt.Errorf("Could not identify type of field: %s", field.Name)
        }
      }
    }
  }

  // Add to list of hash ids
  if _, err = store.Do("RPUSH", pkey, id); err != nil {
    return err
  }

  return nil
}
