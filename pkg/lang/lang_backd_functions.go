package lang

import "github.com/backd-io/backd/backd"

func (l *Lang) addBackdFunctionsCommands() {

	l.AddCommand(
		"obj_create",
		"Creates a data object inside a collection",
		`
Creates a data object inside a collection. Returns ID and error if any.
			
Example:

obj = {}
obj.name = "applicationExample"
obj.description = "this is an example"

id, err = obj_create("myCollection", obj)


	`,
		l.objCreate)

	l.AddCommand(
		"obj_get",
		"Get a data object by ID from a collection",
		`
Get a data object by ID from a collection. Returns object map and error if any.
			
Example:

obj, err = obj_get("myCollection", "idObject")

println(obj.name)

// "applicationExample"
		`,
		l.objGet)

	l.AddCommand(
		"obj_get_many",
		"Get objects from a collection based on the query options passed",
		`
Get objects from a collection based on the query options passed. Returns object array and error if any.
					
Example:

objs, err = obj_get_many("myCollection", { "test": true }, [], 1, 20)

for obj in objs {
	println(pretty(obj))
}
			`,
		l.objGetMany)

	l.AddCommand(
		"obj_all",
		"Get all objects from a collection",
		`
Get all objects from a collection. Returns object array and error if any.
					
Example:

objs, err = obj_all("myCollection")

for obj in objs {
	println(pretty(obj))
}
				`,
		l.objAll)

	l.AddCommand(
		"obj_update",
		"Updates an object on a collection",
		`
Updates an object on a collection. Returns object map and error if any.
						
Example:

obj.name = "newName"

newObj, err = obj_update("myCollection", "idObject", obj)

println(newObj)
		`,
		l.objUpdate)

	l.AddCommand(
		"obj_delete",
		"Deletes an object from a collection",
		`
Deletes an object from a collection. Returns error if any.
						
Example:

err = obj_delete("myCollection", "idObject")

println(err)
		`,
		l.objDelete)

}

// objCreate - obj_create (object)
func (l *Lang) objCreate(collection string, data map[string]interface{}) (string, error) {
	if l.currentAppID == noAppID {
		return "", ErrApplicationNotEspecified
	}
	return l.b.Objects(l.currentAppID).Insert(collection, data)
}

// objGet - obj_get (object)
func (l *Lang) objGet(collection, id string) (data map[string]interface{}, err error) {
	if l.currentAppID == noAppID {
		err = ErrApplicationNotEspecified
		return
	}
	err = l.b.Objects(l.currentAppID).GetByID(collection, id, &data)
	return
}

// objAll - obj_all(collection)
func (l *Lang) objAll(collection string) (data []map[string]interface{}, err error) {
	return l.objGetMany(collection, map[string]interface{}{}, []string{}, 0, 0)
}

// objGetMany - obj_get_many(collection, query, sort, page, per_page)
func (l *Lang) objGetMany(collection string, query map[string]interface{}, sort []string, page, perPage int) (data []map[string]interface{}, err error) {
	if l.currentAppID == noAppID {
		err = ErrApplicationNotEspecified
		return
	}

	var queryOptions backd.QueryOptions
	queryOptions.Q = query
	queryOptions.Sort = sort
	queryOptions.Page = page
	queryOptions.PerPage = perPage

	err = l.b.Objects(l.currentAppID).GetMany(collection, queryOptions, &data)
	return
}

// objUpdate - obj_update (object)
func (l *Lang) objUpdate(collection, id string, data map[string]interface{}) (newData map[string]interface{}, err error) {
	if l.currentAppID == noAppID {
		err = ErrApplicationNotEspecified
		return
	}
	err = l.b.Objects(l.currentAppID).Update(collection, id, data, &newData)
	return
}

// objDelete - obj_delete (object)
func (l *Lang) objDelete(collection, id string) (err error) {
	if l.currentAppID == noAppID {
		err = ErrApplicationNotEspecified
		return
	}
	err = l.b.Objects(l.currentAppID).Delete(collection, id)
	return
}
