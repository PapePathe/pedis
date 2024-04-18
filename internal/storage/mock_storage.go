package storage 

type MockStorage struct {
  GetUserFn func(string) (*User, error)
  DelUserFn func(string) (error)
}

func(ms MockStorage) GetUser(key string)(*User, error) {
  if ms.GetUserFn == nil {
    panic("you must provide an implementation of getUserFn")
  }
  return ms.GetUserFn(key)
}

func(ms MockStorage) SetUser(key string, acl []AclRule)(error) {
  return nil
}

func(ms MockStorage) Users()[]string {
  return nil
}

func(ms MockStorage) DelUser(key string)(error) {
  if ms.DelUserFn == nil {
    panic("you must provide an implementation of delUserFn")
  }
  return ms.DelUserFn(key)
}

func(ms MockStorage) Del(key string)(error) {
  return nil
}

func(ms MockStorage) Get(key string)(string, error) {
  return "", nil
}

func(ms MockStorage) Set(key string, v string, i int64)(error) {
  return nil
}

func(ms MockStorage) HGet(key string)([]byte, error) {
  return nil, nil
}

func(ms MockStorage) HSet(key string, data[]byte, i int64)(int, error) {
  return 0, nil
}

func(ms MockStorage) WriteFromRaft(StorageData)(error) {
  return nil
}
