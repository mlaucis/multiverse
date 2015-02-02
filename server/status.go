/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import "fmt"

// Defining status codes
const (
	// 2XX - Success
	http200OK       = "The request was successful and the response body contains the representation requested."
	http201Created  = "The request was successful, we updated the resource and the response body contains the representation."
	http204NoConent = "The request was successful; the resource was deleted."

	// 3XX - Redirection
	http302Found       = "A common redirect response; you can GET the representation at the URI in the Location response header."
	http304NotModified = "Your client's cached version of the representation is still up to date."

	// 4XX - Client error
	http400BadRequest       = "The data given in the request failed validation. Inspect the response body for details."
	http402Unauthorized     = "The supplied credentials, if any, are not sufficient to access the resource."
	http404NotFound         = "The requested resource was not found on the server"
	http405MethodNotAllowed = "You can't perform this method to the resource."
	http429TooManyRequests  = "Your application is sending too many simultaneous requests."

	// 5XX - Server error
	http500ServerError        = "We couldn't return the representation due to an internal server error."
	http503ServiceUnavailable = "We are temporarily unable to return the representation. Please wait for a bit and try again."
)

// HTTP200OK gets http status
func HTTP200OK() string {
	return fmt.Sprintf(http200OK)
}

// HTTP201Created gets http status
func HTTP201Created() string {
	return fmt.Sprintf(http201Created)
}

// HTTP204NoConent gets http status
func HTTP204NoConent() string {
	return fmt.Sprintf(http204NoConent)
}

// HTTP302Found gets http status
func HTTP302Found() string {
	return fmt.Sprintf(http302Found)
}

// HTTP304NotModified gets http status
func HTTP304NotModified() string {
	return fmt.Sprintf(http304NotModified)
}

// HTTP400BadRequest gets http status
func HTTP400BadRequest() string {
	return fmt.Sprintf(http400BadRequest)
}

// HTTP402Unauthorized gets http status
func HTTP402Unauthorized() string {
	return fmt.Sprintf(http402Unauthorized)
}

// HTTP404NotFound gets http status
func HTTP404NotFound() string {
	return fmt.Sprintf(http404NotFound)
}

// HTTP405MethodNotAllowed gets http status
func HTTP405MethodNotAllowed() string {
	return fmt.Sprintf(http405MethodNotAllowed)
}

// HTTP429TooManyRequests gets http status
func HTTP429TooManyRequests() string {
	return fmt.Sprintf(http429TooManyRequests)
}

// HTTP500ServerError gets http status
func HTTP500ServerError() string {
	return fmt.Sprintf(http500ServerError)
}

// HTTP503ServiceUnavailable gets http status
func HTTP503ServiceUnavailable() string {
	return fmt.Sprintf(http503ServiceUnavailable)
}
