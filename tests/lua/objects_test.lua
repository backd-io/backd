
appid("_backd")

objects = require('backd.objects')

assert(objects.name == "objects")
assert(type(objects) == 'table')
assert(type(objects.new) == 'function')
assert(type(objects.get_one) == 'function')
assert(type(objects.get_many) == 'function')
assert(type(objects.create) == 'function')
assert(type(objects.update) == 'function')
assert(type(objects.delete) == 'function')

col = "collection1"

function printArr(arr)
  for i = 1, table.getn(arr) do
    print(arr[i])
  end
end

-- pre - cleanup
query = {}
items, count = objects.get_many(col, query, {}, 1, 5000)

for i = 1, count do 
  print(string.format("delete.id: '%s'", items[i]._id))
  assert(objects.delete(col, items[i]._id) == true)
end

-- create 5 objects
for i = 1,5 do

  n = objects.new()
  n.name = string.format("John %d", i)
  n.surname = "Doe"

  this = objects.create(col, n)

end

-- test complex query
query = {}
query.name = {}
query.name["$in"] = { "John 1", "John 2" }

items, count = objects.get_many(col, query, {}, 1, 5000)
assert(count == 2)

-- test update
for i = 1, count do
  items[i].name = string.format("Johnny %d", i)
  
  newItem = objects.update(col, items[i]._id, items[i])
  assert(newItem._id == items[i]._id)
  assert(newItem.name == items[i].name)

  -- test get_one
  getItem = objects.get_one(col, items[i]._id)
  assert(getItem._id == items[i]._id)
  assert(getItem.name == items[i].name)
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
