appid("_backd")

json = require('json')
objects = require('backd.objects')
relations = require('backd.relations')

assert(relations.name == "relations")
assert(type(relations) == 'table')
assert(type(relations.get_one_relation) == 'function')
assert(type(relations.get_related) == 'function')
assert(type(relations.get_many_relations) == 'function')
assert(type(relations.create) == 'function')
assert(type(relations.delete) == 'function')

col = "testCollection"

-- helpers
function isInArrByKey(value, items, key) 
  for i = 1, table.getn(items) do
    -- print(items[i][key], value, items[i][key] == value)
    if items[i][key] == value then
      return true
    end
  end
  return false
end

function printEncoded(item)
  print(json.encode(item))
end

-- pre - cleanup
query = {}
items, count = objects.get_many(col, query, {}, 1, 5000)

for i = 1, count do 
  print(string.format("delete.id: '%s'", items[i]._id))
  assert(objects.delete(col, items[i]._id) == true)
end

-- create 5 objects
for i = 1, 50 do
  n = objects.new()
  n.name = string.format("John %d", i)
  n.surname = "Doe"

  this = objects.create(col, n)
end

  -- query items
query = {}
items, count = objects.get_many(col, query, {}, 1, 50)

for i = 2, count do 
  r = objects.new()
  r.src = col 
  r.sid = items[1]._id
  r.dst = col
  r.did = items[i]._id
  r.rel = "link"
  rel = relations.create(r)
  -- printEncoded(rel)
end

related, rel_count = relations.get_many_relations(col, items[1]._id, "out")
assert(rel_count == count - 1)

for o = 2, table.getn(items) do 
  -- ensure the link exists
  assert(isInArrByKey(items[o]._id, related, "did"))
end

linked, linked_count = relations.get_related(col, items[1]._id, "link", "out")
assert(linked_count == count - 1)

for o = 2, table.getn(items) do 
  -- ensure the linked object exists
  assert(isInArrByKey(items[o]._id, linked, "_id"))
  assert(isInArrByKey(items[o].name, linked, "name"))
end

-- assert delete items works
for i = 1, rel_count do 
  assert(relations.delete(related[i]._id) == true)
end

-- post - cleanup
query = {}
items, count = objects.get_many(col, query, {}, 1, 5000)

for i = 1, count do 
  assert(objects.delete(col, items[i]._id), true)
end

-- post - ensure cleanup
query = {}
items, count = objects.get_many(col, query, {}, 1, 5000)
assert(count == 0)

print("tests ok")
