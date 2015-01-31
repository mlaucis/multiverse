#!/bin/bash
curl -i -H "Content-Type: application/json" -d '{"name":"New Account"}' localhost:8082/accounts
curl -i -H "Content-Type: application/json" -d '{"user_name":"User name", "password":"hmac(256)", "email":"de@m.o"}' localhost:8082/account/1/users
curl -i -H "Content-Type: application/json" -d '{"key": "hmac(256)", "name":"New App"}' localhost:8082/account/1/applications
curl -i -H "Content-Type: application/json" -d '{"auth_token": "token1flo", "username": "flo", "name": "Florin", "password": "passwd", "email": "fl@r.in", "url": "blogger", "metadata": "{}"}' localhost:8082/application/1/users
curl -i -H "Content-Type: application/json" -d '{"auth_token": "token1flo", "username": "flo", "name": "Florin", "password": "passwd", "email": "fl@r.in", "url": "blogger", "metadata": "{}"}' localhost:8082/application/1/users
curl -i -H "Content-Type: application/json" -d '{"auth_token": "token1flo", "username": "flo", "name": "Florin", "password": "passwd", "email": "fl@r.in", "url": "blogger", "metadata": "{}"}' localhost:8082/application/1/users
curl -i -H "Content-Type: application/json" -d '{"auth_token": "token1flo", "username": "flo", "name": "Florin", "password": "passwd", "email": "fl@r.in", "url": "blogger", "metadata": "{}"}' localhost:8082/application/1/users
curl -i -H "Content-Type: application/json" -d '{"auth_token": "token1flo", "username": "flo", "name": "Florin", "password": "passwd", "email": "fl@r.in", "url": "blogger", "metadata": "{}"}' localhost:8082/application/1/users
curl -i -H "Content-Type: application/json" -d '{"user_from_id":1,"user_to_id":2}' localhost:8082/application/1/connections
curl -i -H "Content-Type: application/json" -d '{"user_from_id":1,"user_to_id":3}' localhost:8082/application/1/connections
curl -i -H "Content-Type: application/json" -d '{"user_from_id":1,"user_to_id":5}' localhost:8082/application/1/connections
curl -i -H "Content-Type: application/json" -d '{"user_from_id":2,"user_to_id":5}' localhost:8082/application/1/connections
curl -i -H "Content-Type: application/json" -d '{"user_from_id":2,"user_to_id":3}' localhost:8082/application/1/connections
curl -i -H "Content-Type: application/json" -d '{"user_from_id":3,"user_to_id":5}' localhost:8082/application/1/connections
curl -i -H "Content-Type: application/json" -d '{"user_from_id":4,"user_to_id":5}' localhost:8082/application/1/connections
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/application/1/user/3/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/application/1/user/3/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/application/1/user/3/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/application/1/user/5/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/application/1/user/5/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/application/1/user/5/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/application/1/user/5/events
curl -i -H "Content-Type: application/json" -d '{"verb": "like", "metadata": "{}"}' localhost:8082/application/1/user/5/events