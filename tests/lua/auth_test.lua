auth = require('backd.auth')
json = require('json')

assert(auth.name == "auth")
assert(type(auth) == 'table')
assert(type(auth.me) == 'function')
assert(type(auth.login) == 'function')
assert(type(auth.logout) == 'function')
assert(type(auth.get_session_id) == 'function')
assert(type(auth.get_session_state) == 'function')
assert(type(auth.get_session_expiration) == 'function')
assert(type(auth.set_session) == 'function')
assert(type(auth.set_session_id) == 'function')

user = auth.me()
assert(type(user) == 'table')

-- remove - this
encoded = json.encode(user)
print(encoded)



