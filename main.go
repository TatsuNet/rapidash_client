package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"go.knocknote.io/rapidash"
)

type UserLogin struct {
	ID   int64
	Name string
}

func (u *UserLogin) EncodeRapidash(enc rapidash.Encoder) error {
	if u.ID != 0 {
		enc.Int64("id", u.ID)
	}
	enc.String("name", u.Name)
	return enc.Error()
}

func (u *UserLogin) DecodeRapidash(dec rapidash.Decoder) error {
	u.ID = dec.Int64("id")
	u.Name = dec.String("name")
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

// Map column of `user_logins` table to Go type
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
	fmt.Println("Start")

	cache, err := rapidash.New(
		rapidash.ServerAddrs([]string{"localhost:11211"}),
		rapidash.Timeout(3*time.Second),
	)
	if err != nil {
		panic(err)
	}

	tx, err := cache.Begin()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tx.Commit(); err != nil {
			panic(err)
		}
	}()

	userLogin := UserLogin{
		ID:   1,
		Name: "bob",
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
	if err := tx.Find("user_login_slice",  rapidash.Structs(&newUserLoginSlice, new(UserLogin).Struct())); err != nil {
		panic(err)
	}
	fmt.Printf("newUserLoginSlice[0]:%v\n", newUserLoginSlice[0].Name)

}
