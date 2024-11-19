# install `curl`!

# где запущено приложение
export endpoint="http://localhost:3000/"

# эти переменные можно менять
export user_login="alice"
export user_good_password="secret" # успешная авторизация по команде make integraion/auth_ok
export user_wrong_password="not_secret" # ошибка авторизации по команде make integraion/auth_fail

# после успешного логина по команде make integraion/auth_ok
# задайте эту переменную на полученный токен
export session_good_token="88b93361a38167f95a5905ce88e1cd24"

# тут должен быть токен, которого нет в таблице сессий
export session_bad_token="dd5ff7b78319528ad6a8e4d96eed2151"


# данные для создания нового объекта по команде integration/create_ok
# убедитесь, что `session_good_token` задан верно!
export asset_key=key_$(shell date "+%S")
export asset_body=body_$(shell date "+%S")

# после того, как вы создадите объект по команде integration/create_ok, в эту переменную надо сохранить
# идентификатор созданного объекта
export asset_key_good="key_42"

# аутентификация
integration/auth_ok:
	curl -v --data '{"login":$(user_login),"password":$(user_good_password)}' $(endpoint)api/auth

integration/auth_fail:
	curl -v --data '{"login":$(user_login),"password":$(user_wrong_password)}' $(endpoint)api/auth

# попытка создания объекта
integration/create_ok:
	curl -v -H "Authorization: Bearer $(session_good_token)" --data $(asset_body) $(endpoint)api/upload-asset/$(asset_key)

integration/create_fail:
	curl -v -H "Authorization: Bearer $(session_bad_token)" --data $(asset_body) $(endpoint)api/upload-asset/$(asset_key)

# получение объекта
integration/get_ok:
	curl -v -H "Authorization: Bearer $(session_good_token)" $(endpoint)api/asset/$(asset_key_good)

integration/get_fail:
	curl -v -H "Authorization: Bearer $(session_good_token)" $(endpoint)api/asset/not_found
