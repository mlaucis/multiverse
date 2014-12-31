
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
- [x] Implement queries to read/write entities 
- [ ] Implement tests
- [ ] Secure API (Authentication) (check: http://jwt.io/)
- [ ] Add message queue
- [ ] Add cache (Redis)
- [ ] Refine version control (v1, v2 etc.)

# Cases

| Entity | Case | Method | URL | Implementation | Test | Docs |
| ------ | ---- | ------ | --- |:--------------:|:----:|:----:|
|Account|Read Account|GET|`/account/:AccountID`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/account.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Account#read-account)|
|Account|Create account|POST|`/account`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/account.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Account#create-account)|
|Account users|Read account user|GET|`/account/:AccountID/user/:UserID`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/account_user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_user_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Account-users#read-account-user)|
|Account users|Create account user|POST|`/account/:AccountID/user`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/account_user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_user_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Account-users#create-account-user)|
|Account users|List account users|GET|`/account/:AccountID/users`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/account_user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/account_user_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Account-users#list-account-users)|
|Applications|Read application|GET|`/app/:AppID`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/application.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/application_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Applications#read-application)|
|Applications|Create application|POST|`/account/:AccountID/app`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/application.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/application_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Applications#create-application)|
|Applications|List applications|GET|`/account/:AccountID/applications`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/application.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/application_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Applications#list-applications)|
|Users|Read app user|GET|`/app/:AppID/user/:Token`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/user_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Users#read-app-user)|
|Users|Create app user|POST|`/app/:AppID/user`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/user_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Users#create-app-user)|
|Users|List app users|GET|`/app/:AppID/users`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/user.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/user_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Users#list-app-users)|
|User connections|Create user connection|POST|`/app/:AppID/user/:Token/connections`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/connections.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/connections_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/User-connections#create-user-connection)|
|User connections|List user connections|GET|`/app/:AppID/connection`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/connections.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/connections_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/User-connections#list-user-connections)|
|Sessions|Read session|GET|`/app/:AppID/user/:Token/session/:SessionID`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/session.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/session_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Sessions#read-session)|
|Sessions|Create session|POST|`/app/:AppID/user/:userToken/session`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/session.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/session_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Sessions#create-session)|
|Sessions|List user sessions|GET|`/app/:AppID/user/:userToken/sessions`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/session.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/session_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Sessions#list-user-sessions)|
|Events|Read event|GET|`/app/:AppID/event/:EventID`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/event.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Events#read-event)|
|Events|Create event|POST|`/app/:AppID/user/:userToken/session/:SessionID/event`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/event.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Events#create-event)|
|Events|List user events|GET|`/app/:AppID/user/:Token/events`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/event.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Events#list-user-events)|
|Events|List session events|GET|`/app/:AppID/user/:userToken/session/:SessionID/events`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/event.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Events#list-session-events)|
|Events|List connections events|GET|`/app/:AppID/user/:Token/connections/events`|[:white_check_mark:](https://github.com/Tapglue/backend/blob/master/server/event.go)|[:x:](https://github.com/Tapglue/backend/blob/master/server/event_test.go)|[:x:](https://github.com/Tapglue/backend/wiki/Events#list-connections-events)|
