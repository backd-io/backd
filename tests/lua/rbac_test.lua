appid("_backd")

rbac = require('backd.rbac')

assert(rbac.name == "rbac")
assert(type(rbac) == 'table')
assert(type(rbac.get) == 'function')
assert(type(rbac.set) == 'function')
assert(type(rbac.add) == 'function')
assert(type(rbac.remove) == 'function')

