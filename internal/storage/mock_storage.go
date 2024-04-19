package storage 

type MockStorage struct {
  GetUserFn func(string) (*User, error)
  DelUserFn func(string) (error)
  SetUserFn func(string, []AclRule) (error)
  UsersFn func() []string 
  GetFn func(string) (string, error)
}

func(ms MockStorage) GetUser(key string)(*User, error) {
  if ms.GetUserFn == nil {
    panic("you must provide an implementation of getUserFn")
  }
  return ms.GetUserFn(key)
}

func(ms MockStorage) SetUser(key string, acl []AclRule)(error) {
  if ms.SetUserFn == nil {
    panic("you must provide an implementation of SetUserFn")
  }
  return ms.SetUserFn(key, acl)
}

func(ms MockStorage) Users()[]string {
  if ms.UsersFn== nil {
    panic("you must provide an implementation of UsersFn")
  }
  return ms.UsersFn()
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
  if ms.GetFn == nil {
    panic("you must provide an implementation of GetFn")
  }
  return ms.GetFn(key)
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

