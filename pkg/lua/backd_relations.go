package lua

import (
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

// module objects
func (l *Lang) backdRelationsModule(L *lua.LState) int {

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"get_related":        l.getRelated,
		"get_many_relations": l.getManyRelations,
		"get_one_relation":   l.getOneRelation,
		"create":             l.createRelation,
		"delete":             l.deleteRelation,
	})

	L.SetField(mod, "name", lua.LString("relations"))

	// returns the module
	L.Push(mod)
	return 1

}

func (l *Lang) getRelated(L *lua.LState) int {

	var (
		col       string
		id        string
		relation  string
		direction string
		data      map[string]interface{}
		err       error
	)

	col = L.ToString(1)
	id = L.ToString(2)
	relation = L.ToString(3)
	direction = L.ToString(4)

	err = l.b.Objects(l.currentAppID).GetRelationsOf(col, id, relation, direction, &data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(luajson.DecodeValue(l.env, data))
	L.Push(lua.LNumber(len(data)))
	return 2

}

func (l *Lang) getManyRelations(L *lua.LState) int {

	var (
		col       string
		id        string
		direction string
		data      map[string]interface{}
		err       error
	)

	col = L.ToString(1)
	id = L.ToString(2)
	direction = L.ToString(3)

	err = l.b.Objects(l.currentAppID).RelationGetMany(col, id, direction, &data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(luajson.DecodeValue(l.env, data))
	L.Push(lua.LNumber(len(data)))
	return 2

}

func (l *Lang) getOneRelation(L *lua.LState) int {

	var (
		id   string
		data map[string]interface{}
		err  error
	)

	id = L.ToString(1)

	err = l.b.Objects(l.currentAppID).RelationGetByID(id, &data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(luajson.DecodeValue(l.env, data))
	return 1

}

func (l *Lang) createRelation(L *lua.LState) int {

	var (
		data map[string]interface{}
		from []byte
		err  error
	)

	obj := L.CheckAny(1)

	from, err = luajson.Encode(obj)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	data, err = l.b.Objects(l.currentAppID).RelationInsert(from)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(luajson.DecodeValue(l.env, data))
	return 1

}

func (l *Lang) deleteRelation(L *lua.LState) int {

	var (
		id  string
		err error
	)

	id = L.ToString(1)

	err = l.b.Objects(l.currentAppID).RelationDelete(id)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)
	return 1

}
