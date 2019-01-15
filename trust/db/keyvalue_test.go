package db

import (
	"fmt"
	"github.com/go-pg/pg"
	"testing"
)

func TestKeyValue(t *testing.T) {
	DatabaseUrl := fmt.Sprintf("postgres://postgres:@localhost/crosschain_%d", 1)
	rDb := &DB{}
	info, err := pg.ParseURL(DatabaseUrl)
	if err != nil {
		t.Error(err)
	}
	println(info.Network, info.Addr)
	rDb.Connect(info.Addr, info.User, info.Password, info.Database)
	rDb.Debug()

	err = UpdateValue(rDb, "key", "val1")
	if err != nil {
		t.Error(err)
	}

	v := GetValue(rDb, "key")
	if v == nil {
		t.Error("should found value")
	}

	println(v.Id, v.Value)

	err = UpdateValue(rDb, "key", "val2")
	v = GetValue(rDb, "key")
	if v == nil {
		t.Error("should found value")
	}
	if v.Value != "val2" {
		t.Error("Unexpected value")
	}

	println(v.Id, v.Value)
}
