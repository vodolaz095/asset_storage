# install `curl`!

# эти переменные можно менять
export user_login="alice"
export user_good_password="secret"
export user_wrong_password="not_secret"

# после успешного логина задайте эту переменную на полученный, в результате вызова
# `make integraion/auth_ok`, токен
export session_good_token="dd5f8e4d96eed2151f7b78319528ad6a"

# тут должен быть токен, которого нет в таблице сессий
export session_bad_token="dd5ff7b78319528ad6a8e4d96eed2151"

# данные для создания нового объекта
export asset_key=key_$(shell date "+%S")
export asset_body=body_$(shell date "+%S")

export asset_key_good="key_01"

integration/auth_ok:
	curl -v --data '{"login":$(user_login),"password":$(user_good_password)}' http://localhost:3000/api/auth

integration/auth_fail:
	curl -v --data '{"login":$(user_login),"password":$(user_wrong_password)}' http://localhost:3000/api/auth

integration/create_ok:
	curl -v -H "Authorization: Bearer $(session_good_token)" --data $(asset_body) http://localhost:3000/api/upload-asset/$(asset_key)

integration/create_fail:
	curl -v -H "Authorization: Bearer $(session_bad_token)" --data $(asset_body) http://localhost:3000/api/upload-asset/$(asset_key)

integration/get_ok:
	curl -v -H "Authorization: Bearer $(session_good_token)"  http://localhost:3000/api/asset/$(asset_key_good)

integration/get_fail:
	curl -v -H "Authorization: Bearer $(session_good_token)"  http://localhost:3000/api/asset/not_found
