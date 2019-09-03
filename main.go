package main

import (
	"fmt"
	"github.com/pkg/errors"
	"time"

	"go.knocknote.io/rapidash"
)

type UserLogin struct {
	ID            int64
	UserID        int64
	UserSessionID int64
	LoginParamID  int64
	Name          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (u *UserLogin) EncodeRapidash(enc rapidash.Encoder) error {
	if u.ID != 0 {
		enc.Int64("id", u.ID)
	}
	enc.Int64("user_id", u.UserID)
	enc.Int64("user_session_id", u.UserSessionID)
	enc.Int64("login_param_id", u.LoginParamID)
	enc.String("name", u.Name)
	enc.Time("created_at", u.CreatedAt)
	enc.Time("updated_at", u.UpdatedAt)
	return enc.Error()
}

func (u *UserLogin) DecodeRapidash(dec rapidash.Decoder) error {
	u.ID = dec.Int64("id")
	u.UserID = dec.Int64("user_id")
	u.UserSessionID = dec.Int64("user_session_id")
	u.LoginParamID = dec.Int64("login_param_id")
	u.Name = dec.String("name")
	u.CreatedAt = dec.Time("created_at")
	u.UpdatedAt = dec.Time("updated_at")
	return dec.Error()
}

type UserLoginSlice []*UserLogin

func (u *UserLoginSlice) EncodeRapidash(enc rapidash.Encoder) error {
	for _, v := range *u {
		if err := v.EncodeRapidash(enc.New()); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (u *UserLoginSlice) DecodeRapidash(dec rapidash.Decoder) error {
	decLen := dec.Len()
	*u = make([]*UserLogin, decLen)
	for i := 0; i < decLen; i++ {
		var v UserLogin
		if err := v.DecodeRapidash(dec.At(i)); err != nil {
			return errors.WithStack(err)
		}
		(*u)[i] = &v
	}
	return nil
}

func (u *UserLogin) Struct() *rapidash.Struct {
	return rapidash.NewStruct("user_logins").
		FieldInt64("id").
		FieldInt64("user_id").
		FieldInt64("user_session_id").
		FieldInt64("login_param_id").
		FieldString("name").
		FieldTime("created_at").
		FieldTime("updated_at")
}

func main() {
	cache, err := rapidash.New(
		rapidash.ServerAddrs([]string{"localhost:11211"}),
		rapidash.Timeout(3*time.Second),
	)
	if err != nil {
		panic(err)
	}

	tx, err := cache.Begin()
	defer func() {
		if err := tx.Commit(); err != nil {
			fmt.Printf("err:%v\n", err)
		}
	}()
	if err != nil {
		panic(err)
	}

	if err := tx.Create("key", rapidash.Int(1)); err != nil {
		panic(err)
	}

	var v int
	if err := tx.Find("key", rapidash.IntPtr(&v)); err != nil {
		panic(err)
	}
	fmt.Printf("v:%v\n", v)

	userLogin := UserLogin{
		ID:     1,
		UserID: 1,
		Name:   "bob",
	}

	if err := tx.Create("user_login", new(UserLogin).Struct().Cast(&userLogin)); err != nil {
		panic(err)
	}
	var newUserLogin UserLogin
	if err := tx.Find("user_login", new(UserLogin).Struct().Cast(&newUserLogin)); err != nil {
		panic(err)
	}
	fmt.Printf("newUserLogin:%v\n", newUserLogin.Name)

	var userLoginSlice UserLoginSlice
	userLoginSlice = append(userLoginSlice, &userLogin)

	if err := tx.Create("user_login_slice", rapidash.Structs(&userLoginSlice, new(UserLogin).Struct())); err != nil {
		panic(err)
	}
	var newUserLoginSlice UserLoginSlice
	if err := tx.Find("user_login_slice", rapidash.Structs(&newUserLoginSlice, new(UserLogin).Struct())); err != nil {
		panic(err)
	}
	fmt.Printf("newUserLoginSlice[0]:%v\n", newUserLoginSlice[0].Name)

}
