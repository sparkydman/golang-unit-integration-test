package error_formats

import (
	"efficient_api/utils/error_utils"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func ParseError(err error) error_utils.MessageErr{
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), "no rows in result set"){
			return error_utils.NewNotFoundErrorMessage("No rows found matching the id")
		}
		return error_utils.NewInternalServerErrorMessage(fmt.Sprintf("Error occured rying to proess request: %s", err.Error()))
	}

	switch sqlErr.Number{
	case 1062:
		return error_utils.NewBadRequestErrorMessage("Title already taken")
	}
	return error_utils.NewInternalServerErrorMessage(fmt.Sprintf("Error occured rying to proess request: %s", err.Error()))
}