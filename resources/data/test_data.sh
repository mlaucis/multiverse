#!/bin/bash
curl -i -H "Content-Type: application/json" -d '{"name":"New Account", "description":"Description of the account"}' localhost:8082/0.1/accounts
curl -i -H "Content-Type: application/json" -H "Authorization: Bearer token_1_TmV3IEFjY291bnQ=" -d '{"user_name":"User name", "first_name": "Demo", "last_name": "User", "password":"hmac(256)", "email":"de@m.o"}' localhost:8082/0.1/account/1/users
curl -i -H "Content-Type: application/json" -d '{"key": "hmac(256)", "name":"New App","description":"awesomeness"}' localhost:8082/0.1/account/1/applications
curl -i -H "Content-Type: application/json" -d '{"auth_token": "yZg6ZCJjHGy5caTcVnD25pVMEswUEQWTSA64tkBU", "user_name": "dlsniper", "first_name": "Florin", "last_name": "Patan", "password": "JjQxWYYCWfX634q6KQeDSusywVH3T5Dw9hMqBcUd", "email": "florinpatan@gmail.com", "url": "http://florinpatan.ro", "metadata": "{}"}' localhost:8082/0.1/application/1/users
curl -i -H "Content-Type: application/json" -d '{"auth_token": "yZg6ZCJjHGy5caTcVnD25pVMEswUEQWTSA64tkBU", "user_name": "dlsniper", "first_name": "Florin", "last_name": "Patan", "password": "JjQxWYYCWfX634q6KQeDSusywVH3T5Dw9hMqBcUd", "email": "florinpatan@gmail.com", "url": "http://florinpatan.ro", "metadata": "{}"}' localhost:8082/0.1/application/1/users
curl -i -H "Content-Type: application/json" -d '{"auth_token": "yZg6ZCJjHGy5caTcVnD25pVMEswUEQWTSA64tkBU", "user_name": "dlsniper", "first_name": "Florin", "last_name": "Patan", "password": "JjQxWYYCWfX634q6KQeDSusywVH3T5Dw9hMqBcUd", "email": "florinpatan@gmail.com", "url": "http://florinpatan.ro", "metadata": "{}"}' localhost:8082/0.1/application/1/users
curl -i -H "Content-Type: application/json" -d '{"auth_token": "yZg6ZCJjHGy5caTcVnD25pVMEswUEQWTSA64tkBU", "user_name": "dlsniper", "first_name": "Florin", "last_name": "Patan", "password": "JjQxWYYCWfX634q6KQeDSusywVH3T5Dw9hMqBcUd", "email": "florinpatan@gmail.com", "url": "http://florinpatan.ro", "metadata": "{}"}' localhost:8082/0.1/application/1/users
curl -i -H "Content-Type: application/json" -d '{"auth_token": "yZg6ZCJjHGy5caTcVnD25pVMEswUEQWTSA64tkBU", "user_name": "dlsniper", "first_name": "Florin", "last_name": "Patan", "password": "JjQxWYYCWfX634q6KQeDSusywVH3T5Dw9hMqBcUd", "email": "florinpatan@gmail.com", "url": "http://florinpatan.ro", "metadata": "{}"}' localhost:8082/0.1/application/1/users
curl -i -H "Content-Type: application/json" -d '{"user_from_id":1,"user_to_id":2}' localhost:8082/0.1/application/1/user/1/connections
curl -i -H "Content-Type: application/json" -d '{"user_from_id":1,"user_to_id":3}' localhost:8082/0.1/application/1/user/1/connections
curl -i -H "Content-Type: application/json" -d '{"user_from_id":1,"user_to_id":5}' localhost:8082/0.1/application/1/user/1/connections
curl -i -H "Content-Type: application/json" -d '{"user_from_id":2,"user_to_id":5}' localhost:8082/0.1/application/1/user/2/connections
curl -i -H "Content-Type: application/json" -d '{"user_from_id":2,"user_to_id":3}' localhost:8082/0.1/application/1/user/2/connections
curl -i -H "Content-Type: application/json" -d '{"user_from_id":3,"user_to_id":5}' localhost:8082/0.1/application/1/user/3/connections
curl -i -H "Content-Type: application/json" -d '{"user_from_id":4,"user_to_id":5}' localhost:8082/0.1/application/1/user/3/connections
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/0.1/application/1/user/1/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/0.1/application/1/user/2/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/0.1/application/1/user/3/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/0.1/application/1/user/4/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/0.1/application/1/user/5/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/0.1/application/1/user/5/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/0.1/application/1/user/5/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/0.1/application/1/user/5/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/0.1/application/1/user/1/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/0.1/application/1/user/1/events
curl -i -H "Content-Type: application/json" -d '{"user_from_id":5,"user_to_id":1}' localhost:8082/0.1/application/1/user/5/connections
curl -i -H "Content-Type: application/json" -d '{"token":"token_1_TmV3IEFjY291bnQ=", "name":"New Account","description":"Another description of the account", "enabled": true, "created_at":"2015-02-02T19:13:18.239759449Z", "received_at":"2015-02-02T19:13:18.239759449Z", "metadata":"{123}"}' -X PUT localhost:8082/0.1/account/1
curl -i -H "Content-Type: application/json" -d '{"name":"New Account", "description":"Description of the account"}' localhost:8082/0.1/accounts
curl -i -X DELETE localhost:8082/0.1/account/2
curl -i -H "Content-Type: application/json" -H "Authorization: Bearer token_1_TmV3IEFjY291bnQ=" -d '{"user_name":"User name", "first_name": "Demo", "last_name": "User", "password":"hmac(256)", "email":"de@m.o"}' localhost:8082/0.1/account/1/users
curl -i -H "Content-Type: application/json" -H "Authorization: Bearer token_1_TmV3IEFjY291bnQ=" -d '{"user_name":"User name changed", "first_name": "Demo", "last_name": "User", "password":"hmac(256)changed", "email":"de@m.ohno"}' -X PUT localhost:8082/0.1/account/1/user/2
curl -i -X DELETE localhost:8082/0.1/account/1/user/2
curl -i -H "Content-Type: application/json" -d '{"key": "hmac(256)", "name":"New App","description":"awesomeness"}' localhost:8082/0.1/account/1/applications
curl -i -H "Content-Type: application/json" -d '{"key": "hmac(256)", "name":"New App","description":"awesomeness"}' localhost:8082/0.1/account/1/applications
curl -i -H "Content-Type: application/json" -d '{"key": "hmac(256)", "name":"New App changed","description":"awesomeness changed"}' -X PUT localhost:8082/0.1/account/1/application/2
curl -i -X DELETE localhost:8082/0.1/account/1/application/2