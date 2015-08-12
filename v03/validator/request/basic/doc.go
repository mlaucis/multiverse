// Package basic checks if the request is signed with the proper credentials using HTTP Basic Authentication.
//
// For an account / account user / application request, the request will be signed
// as follows:
//  - request has an account context -> Authorization: Basic base64(accountKey)
//  - request has an account + account user context -> Authorization: Basic base64(accountKey:accountUser)
//  - request has an application context -> Authorization: Basic base64(applicationKey)
//  - request has an application + application user context -> Authorization: Basic base64(applicationKey:applicationUser)
//
// In order to receive the account and accountUser keys, one must call the account login endpoint.
//
// In order to receive the application key, one must create a new application in the dashboard, which must be
// stored in a highly secure fashion.
//
// In order to receive the applicationUser key, one must call the application login endpoint.
package basic
