package common

import (
	"github.com/lib/pq"
)

func AddUniqueIdToArray(array pq.Int64Array, newId int64) pq.Int64Array {
	for _, id := range array {
		if id == newId {
			return array
		}
	}

	return append(array, newId)
}
