/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

/*
Package request holds the whole validation routines for the requests.

It supports multiple authentication / request signature methods.

Detection of the used method is done by searching for various headers
used by each of the supported authentication providers and executing the
according checker.

*/
package request
