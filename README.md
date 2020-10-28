![one logo](https://github.com/catmullet/one/blob/assets/one_logo.png)
# _(īdəmˌpōtənt)_ Idempotency Handler
## Description
Easily check if you have recieved/processed an object before. Some examples may be:
* Where PubSub services use _"at least once delivery"_
* Cases of accepting requests to make payment
* Deduping requests to your services

## Getting Started
import the repo by running:
```sh
go get github.com/catmullet/one
```
import into your project
```go
import github.com/catmullet/one
```
## Creating Keys
The `one.MakeKey()` function takes in an array of any type and creates a key. This key is specific to the parameters you have passed in.  If you pass in the exact same fields you will get the exact same key.  Its how we tell if we are getting the same request.  Choose Something that will give you a good indication that this is not the same object.  you can be as strict or relaxed as you want with it. 

Take for example this event from cloud storage
```go
type Event struct {
  Bucket string
  Object string
  Version int
  UpdateTime time.Time
}
```
If you wanted to make sure that you never processed this storage object again you would use this:
```go
key := one.MakeKey(event.Bucket, event.Object)
```
If you wanted to make sure you processed on every version update you would use this:
```go
key := one.MakeKey(event.Bucket, event.Object, event.Version)
```
If you wanted to process based on any change to the storage object you could pass in the entire object like this:
```go
key := one.MakeKey(event)
```
## Implementing the OneStore
Once you have created your key you're going to need a place to store it. This is so you can check against it later on.
Currently there are only two implementations of storage.  
#### Local Storage _(Only for single instance Apps)_
```go
var oneStore OneStore
oneStore = localstore.NewLocalOneStore(ctx, time.Minute * 5)
```
#### Redis
```go
// import "gopkg.in/redis.v5" for redis.Options

options := &redis.Options{
		Network:            "tcp",
		Addr:               "localhost:6379",
		Dialer:             nil,
		Password:           "",
		DB:                 0,
	}
  
var oneStore OneStore
oneStore = redisstore.NewRedisOneStore(options, time.Second * 30)
```
## Add Keys
```go
err := oneStore.Add(key)
if err == one.ErrKeyExists {
  // Key already exists, so handle that here.
}

// Key doesn't exist and was added to the one store
```
## Checking Keys
err := oneStore.HasKey(key)
switch err {
  case 
if err == one.ErrNoKeyExists {
  // Key does not exist in the one store
}

// Key exists
