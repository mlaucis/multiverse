
# To Do's

- [x] Draft API design
- [x] Design entities
- [x] Create db tables
- [x] Implement routes
- [x] Mock API responses
- [x] Document code
- [x] Document api
- [x] Implement basic logging
- [x] Adjust response headers and caching
- [x] Create test files
- [x] ORM + DB Connection
- [x] Implement request validation in routes with regex 
- [ ] Implement queries to read/write entities 
- [ ] Implement tests
- [ ] Secure API (Authentication) (check: http://jwt.io/)
- [ ] Add message queue
- [ ] Add cache (Redis)
- [ ] Refine version control (v1, v2 etc.)

# Cases

| Entity | Case | Method | URL | Implementation | Test | Docs |
| ------ | ---- | ------ | --- |:--------------:|:----:|:----:|
|Account|Read Account|GET|`/account/:AccountID`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/account.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/account_test.go)|[:x:](https://github.com/Gluee/backend/wiki/1.-Account#read-account)|
|Account|Create account|POST|`/account`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/account.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/account_test.go)|[:x:](https://github.com/Gluee/backend/wiki/1.-Account#create-account)|
|AccountUser|Read account user|GET|`/account/:AccountID/user/:UserID`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/account_user.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/account_user_test.go)|[:x:](https://github.com/Gluee/backend/wiki/2.-AccountUser#read-account-user)|
|AccountUser|Create account user|POST|`/account/:AccountID/user`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/account_user.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/account_user_test.go)|[:x:](https://github.com/Gluee/backend/wiki/2.-AccountUser#create-account-user)|
|AccountUser|List account users|GET|`/account/:AccountID/users`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/account_user.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/account_user_test.go)|[:x:](https://github.com/Gluee/backend/wiki/2.-AccountUser#list-account-users)|
|Application|Read application|GET|`/app/:AppID`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/application.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/application_test.go)|[:x:](https://github.com/Gluee/backend/wiki/3.-Application#read-application)|
|Application|Create application|POST|`/account/:AccountID/app`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/application.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/application_test.go)|[:x:](https://github.com/Gluee/backend/wiki/3.-Application#create-application)|
|Application|List applications|GET|`/account/:AccountID/applications`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/application.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/application_test.go)|[:x:](https://github.com/Gluee/backend/wiki/3.-Application#list-applications)|
|User|Read app user|GET|`/app/:AppID/user/:Token`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/user.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/user_test.go)|[:x:](https://github.com/Gluee/backend/wiki/4.-User#read-app-user)|
|User|Create app user|POST|`/app/:AppID/user`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/user.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/user_test.go)|[:x:](https://github.com/Gluee/backend/wiki/4.-User#create-app-user)|
|User|List app users|GET|`/app/:AppID/users`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/user.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/user_test.go)|[:x:](https://github.com/Gluee/backend/wiki/4.-User#list-app-users)|
|UserConnection|Create user connection|POST|`/app/:AppID/user/:Token/connections`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/connections.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/connections_test.go)|[:x:](https://github.com/Gluee/backend/wiki/5.-UserConnection#create-user-connection)|
|UserConnection|List user connections|GET|`/app/:AppID/connection`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/connections.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/connections_test.go)|[:x:](https://github.com/Gluee/backend/wiki/5.-UserConnection#list-user-connections)|
|Session|Read session|GET|`/app/:AppID/user/:Token/session/:SessionID`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/session.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/session_test.go)|[:x:](https://github.com/Gluee/backend/wiki/6.-Session#read-session)|
|Session|Create session|POST|`/app/:AppID/user/:userToken/session`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/session.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/session_test.go)|[:x:](https://github.com/Gluee/backend/wiki/6.-Session#create-session)|
|Session|List user sessions|GET|`/app/:AppID/user/:userToken/sessions`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/session.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/session_test.go)|[:x:](https://github.com/Gluee/backend/wiki/6.-Session#list-user-sessions)|
|Event|Read event|GET|`/app/:AppID/event/:EventID`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/event.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Gluee/backend/wiki/7.-Event#read-event)|
|Event|Create event|POST|`/app/:AppID/user/:userToken/session/:SessionID/event`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/event.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Gluee/backend/wiki/7.-Event#create-event)|
|Event|List user events|GET|`/app/:AppID/user/:Token/events`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/event.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Gluee/backend/wiki/7.-Event#list-user-events)|
|Event|List session events|GET|`/app/:AppID/user/:userToken/session/:SessionID/events`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/event.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Gluee/backend/wiki/7.-Event#list-session-events)|
|Event|List connections events|GET|`/app/:AppID/user/:Token/connections/events`|[:white_check_mark:](https://github.com/Gluee/backend/blob/master/server/event.go)|[:x:](https://github.com/Gluee/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Gluee/backend/wiki/7.-Event#list-connections-events)|
