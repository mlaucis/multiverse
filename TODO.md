
# To Do's

1. ORM + DB Connection
2. Refine API documentation
3. Code documentation

# Use cases

| Entity | Verb | Case | Implementation | Test | Documentation |
| ------ | ---- | ---- | -------------- | ---- | ------------- |
| Account|POST|Create an Account|✓|✗|✗|✗|
| Account|GET|Read an Account|✓|✗|✗|✗|
| AccountUser|POST|Create an Account User|✗|✗|✗|✗|
| AccountUser|GET|List Account Users|✗|✗|✗|✗|
| AccountUser|GET|Read Single Account User|✗|✗|✗|✗|
| Application|POST|Create an Application|✗|✗|✗|✗|
| Application|GET|List Applications of Account|✗|✗|✗|✗|
| Application|GET|Read Single Application|✗|✗|✗|✗|
| User|POST|Create a User|✗|✗|✗|✗|
| User|GET|List Users of App|✗|✗|✗|✗|
| User|GET|Read Single User|✗|✗|✗|✗|
| UserConnection|POST|Create a User Connection|✗|✗|✗|✗|
| UserConnection|GET|List Users Connections|✗|✗|✗|✗|
| Session|POST|Create a Session|✗|✗|✗|✗|
| Session|GET|List Sessions of a User|✗|✗|✗|✗|
| Session|GET|Read Single Session|✗|✗|✗|✗|
| Event|POST|Create an event|✗|✗|✗|✗|
| Event|GET|List events of a user|✗|✗|✗|✗|
| Event|GET|List events of a session|✗|✗|✗|✗|
| Event|GET|Read single event|✗|✗|✗|✗|
| Event|GET|List events of a users connections|✗|✗|✗|✗|

# Pipeline

- Extending components (MQ, Redis, etc.)
- Tests
- Secure API (Authentication) (check: http://jwt.io/)
- Version Control (v1, v2 etc.)