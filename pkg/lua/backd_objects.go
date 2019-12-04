package lua

import (
	"encoding/json"

	"github.com/backd-io/backd/backd"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

// module objects
func (l *Lang) backdObjectsModule(L *lua.LState) int {

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"new":      l.newObject,
		"get_one":  l.getObject,
		"get_many": l.getObjects,
		"create":   l.createObject,
		"update":   l.updateObject,
		"delete":   l.deleteObject,
	})

	L.SetField(mod, "name", lua.LString("objects"))

	// returns the module
	L.Push(mod)
	return 1

}

func (l *Lang) newObject(L *lua.LState) int {
	L.Push(new(lua.LTable))
	return 1
}

func (l *Lang) getObject(L *lua.LState) int {

	var (
		col  string
		id   string
		data map[string]interface{}
		err  error
	)

	col = L.ToString(1)
	id = L.ToString(2)

	err = l.b.Objects(l.currentAppID).GetByID(col, id, &data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(luajson.DecodeValue(l.env, data))
	return 1
}

// get_many(collection, query, sort, page, per_page)
//
// ex:
//
// query = {}
// query.name = "test"
// sort = {}
// items, count = objects.get_many("collection", query, sort, 1, 20)
func (l *Lang) getObjects(L *lua.LState) int {

	var (
		col          string
		b, b2        []byte
		query        map[string]interface{}
		queryOptions backd.QueryOptions
		sort         []string
		data         []interface{}
		err          error
	)

	col = L.ToString(1)

	// query
	b, err = luajson.Encode(L.CheckAny(2))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	err = json.Unmarshal(b, &query)
	if err != nil {
		query = map[string]interface{}{}
	}

	// sort
	b2, err = luajson.Encode(L.CheckAny(3))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	err = json.Unmarshal(b2, &sort)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	queryOptions.Q = query
	queryOptions.Sort = sort
	queryOptions.Page = L.CheckInt(4)
	queryOptions.PerPage = L.CheckInt(5)

	err = l.b.Objects(l.currentAppID).GetMany(col, queryOptions, &data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(luajson.DecodeValue(l.env, data))
	L.Push(lua.LNumber(len(data)))
	return 2

}

func (l *Lang) createObject(L *lua.LState) int {

	var (
		col  string
		from []byte
		to   map[string]interface{}
		err  error
	)

	col = L.ToString(1)
	obj := L.CheckAny(2)

	from, err = luajson.Encode(obj)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	to, err = l.b.Objects(l.currentAppID).Insert(col, from)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(luajson.DecodeValue(l.env, to))
	return 1

}

func (l *Lang) updateObject(L *lua.LState) int {

	var (
		col  string
		id   string
		from []byte
		to   map[string]interface{}
		err  error
	)

	col = L.ToString(1)
	id = L.ToString(2)
	obj := L.CheckAny(3)

	from, err = luajson.Encode(obj)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	err = l.b.Objects(l.currentAppID).Update(col, id, from, &to)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(luajson.DecodeValue(l.env, to))
	return 1

}

func (l *Lang) deleteObject(L *lua.LState) int {

	var (
		col string
		id  string
		err error
	)

	col = L.ToString(1)
	id = L.ToString(2)

	err = l.b.Objects(l.currentAppID).Delete(col, id)

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1

}
