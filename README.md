# pet проект!

## Простое веб api для доступа к elastic search engine

Настраивается через creds.yaml\
`/places` - вывод всей информации с пагинацией по 10 записей, с переходом по страницам
`/api/places` - аналогично первому, формат JSON
`/api/get_token` - получение токена для jwt
`/api/recommend?lat=55.674&lon=37.666` - только с токеном, получение 3 ближайших мест, через указание текущих lon и lat

TODO: 
Header content -> JSON
Find and add error handler