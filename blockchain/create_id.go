package blockchain

import (
	"strconv"

	"svcledger/helpers"
)

func CreateId(pk string, timestamp int64) string {
	return helpers.CreateHash([]byte(pk + ";" + strconv.FormatInt(timestamp, 16)))
}
