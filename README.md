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
	SetValue(v *Value) error // set session data
	GetValue() *Value        // get sesstion data
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

### Features

- [x] Encrypt cookie value with CFB alg
- [x] HTTP Only

### Authors

- [Spring MC](https://twitter.com/mcspring)