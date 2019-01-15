backd = require('backd')
json = require('json')

assert(backd.name == "backd")
assert(type(backd) == 'table')
assert(type(backd.me) == 'function')
-- assert(type(backd.get_one) == 'function')
-- assert(type(backd.get_many) == 'function')
-- assert(type(backd.create) == 'function')
-- assert(type(backd.update) == 'function')
-- assert(type(backd.delete) == 'function')

user = backd.me()
assert(type(user) == 'table')

encoded = json.encode(user)
print(encoded)

function keys2(table)
  arr = {}
  for k,v in pairs(table) do
    table.insert(arr, k)
  end
  return arr
end

