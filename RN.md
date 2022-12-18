curl --location --request POST 'http://localhost:8080/v1/api/history_update' \
--header 'Content-Type: application/json' \
--data-raw '{
"description_ru": "Подключена возможность оплаты подписки  с CryptoCloud",
"description_en": "The possibility of paying for a subscription with CryptoCloud is connected"
}'

curl --location --request POST 'http://localhost:8080/v1/api/history_update' \
--header 'Content-Type: application/json' \
--data-raw '{
"description_ru": "Адаптация меню для мобильных устройств.\r\nТаблица антропометрии адаптирована под мобильное приложение",
"description_en": "Menu adaptation for mobile devices.\n Anthropometry table adapted for mobile application"
}'

curl --location --request POST 'http://localhost:8080/v1/api/history_update' \
--header 'Content-Type: application/json' \
--data-raw '{
"description_ru": "Добавлена история изменений",
"description_en": "Added history of changes"
}'

curl --location --request POST 'https://chartdrug.com/v1/api/history_update' \
--header 'Content-Type: application/json' \
--data-raw '{
"description_ru": "Добавлен адаптивный дизайн для фармакологии",
"description_en": "Added adaptive design for pharmacology"
}'