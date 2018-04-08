# session

[![Build Status](https://travis-ci.org/dolab/session.svg?branch=master)](https://travis-ci.org/dolab/session)

Session manager for golang web application.

### Howto

Session defines an interface allowing user customing there own `Provider` and `Store`.

```go
// Provider defines session store provider apis
type Provider interface {
	New(sid string) (sto Storer, err error)
	Restore(sid string) (sto Storer, err error)
	Refresh(sid, newsid string) (Storer, error)
	Destroy(sid string) error
}

// Storer defines session store apis
type Storer interface {
	SessionID() string       // return current session ID
	SetValue(v *session.Value) error // set session data
	GetValue() *session.Value        // get sesstion data
	Touch() error            // sync session expire time to the provider
}
```

### Usage

```go
func (_ *User) Login(ctx *gogo.Context) {
    sess := session.New(provider, config)

    sto, err := sess.Start(ctx.Response, ctx.Request)
    if err != nil {
        ctx.SetStatush(http.StatusInternalError)
        return
    }

    // do user verification logic

    // save current user to session store
    sto.GetValue().Set("current_user", user)

    ctx.Return()
}
```

- use redis provider 
``` go
import(
    providers "github.com/dolab/session/providers/redis"
    
    "github.com/go-redis/redis"
)

func (_ *User) Login(ctx *gogo.Context) {
    client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

    sess := session.New(providers.New(client), config)

    sto, err := sess.Start(ctx.Response, ctx.Request)
    if err != nil {
        ctx.SetStatush(http.StatusInternalError)
        return
    }

    // do user verification logic

    // save current user to session store
    sto.GetValue().Set("current_user", user)

    ctx.Return()
}
```

- use memcache provider
``` go
import(
    providers "github.com/dolab/session/providers/memcache"

	"github.com/bradfitz/gomemcache/memcache"
)

func (_ *User) Login(ctx *gogo.Context) {
    client := memcache.New("127.0.0.1:11211")

    sess := session.New(providers.New(client), config)

    sto, err := sess.Start(ctx.Response, ctx.Request)
    if err != nil {
        ctx.SetStatush(http.StatusInternalError)
        return
    }

    // do user verification logic

    // save current user to session store
    sto.GetValue().Set("current_user", user)

    ctx.Return()
}
```

### Features

- [x] Encrypt cookie value with CFB alg
- [x] HTTP Only

### Authors

- [Spring MC](https://twitter.com/mcspring)