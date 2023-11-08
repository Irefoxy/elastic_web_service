# pet проект!

## Простое веб api для доступа к elastic search engine

Настраивается через creds.yaml\
`/api/places` - вывод всей информации с пагинацией по 10 записей, формат JSON
`/api/get_token` - получение токена для jwt
`/api/recommend?lat=55.674&lon=37.666` - только с токеном, получение 3 ближайших мест, через указание текущих lon и lat\
Проверка работы в терминале `curl -X GET -H "Authorization: Bearer <токен>" "http://<адрес:порт сервера>/api/recommend?lat=55.674&lon=37.666"`
TODO: 
