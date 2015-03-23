/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package keys

// getScope returns the full-formed request scope and version
func getScope(date, scope, requestVersion string) string {
	return date + "/" + scope + "/" + requestVersion
}
