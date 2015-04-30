/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import "fmt"

func appSchema(query string, accountID, applicationID int64) string {
	return fmt.Sprintf(query, accountID, applicationID)
}
