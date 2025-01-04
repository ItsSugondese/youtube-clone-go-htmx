package generic_repo

import (
	pagination_utils "youtube-clone/pkg/utils/pagination-utils"
	"math"
	"strconv"

	"gorm.io/gorm"
)

func Paginate(value interface{}, pagination *pagination_utils.PaginationRequest, paginationResponse *pagination_utils.PaginationResponse,
	db *gorm.DB, preloadAssociations ...string) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	paginationResponse.TotalElements = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Rows)))
	paginationResponse.TotalPages = totalPages
	paginationResponse.CurrentPageIndex = pagination.Page

	// since total no. of elements actual elements can't be GET without executing query, we're calculating it with assumptions
	convertedRows := int(totalRows)
	totalExpectedElements := pagination.Page * pagination.Rows
	if totalExpectedElements > convertedRows {
		paginationResponse.NoOfElements = convertedRows - (pagination.Page-1)*pagination.Rows
	} else {
		paginationResponse.NoOfElements = pagination.Rows
	}

	//return func(db *gorm.DB) *gorm.DB {
	//	return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	//}

	return func(db *gorm.DB) *gorm.DB {
		// Apply preloads
		for _, association := range preloadAssociations {
			db = db.Preload(association)
		}
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func RawQueryPaginate(pagination *pagination_utils.PaginationRequest, paginationResponse *pagination_utils.PaginationResponse,
	db *gorm.DB, rawSQL string, args ...interface{}) func(db *gorm.DB) *gorm.DB {
	var totalRows int64

	if args != nil {
		countSQL := "SELECT COUNT(*) FROM (" + rawSQL + ") AS subquery"
		db.Raw(countSQL, args...).Scan(&totalRows)
	} else {

		db.Raw(rawSQL).Count(&totalRows)
	}

	paginationResponse.TotalElements = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Rows)))
	paginationResponse.TotalPages = totalPages
	paginationResponse.CurrentPageIndex = pagination.Page

	// since total no. of elements actual elements can't be GET without executing query, we're calculating it with assumptions
	convertedRows := int(totalRows)
	totalExpectedElements := pagination.Page * pagination.Rows
	if totalExpectedElements > convertedRows {
		paginationResponse.NoOfElements = convertedRows - (pagination.Page-1)*pagination.Rows
	} else {
		paginationResponse.NoOfElements = pagination.Rows
	}

	//return func(db *gorm.DB) *gorm.DB {
	//	return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	//}

	return func(db *gorm.DB) *gorm.DB {
		// Apply pagination using raw SQL
		if rawSQL != "" {
			argsLength := len(args)
			limitPlaceholder := "$" + strconv.Itoa(argsLength+1)
			offsetPlaceholder := "$" + strconv.Itoa(argsLength+2)

			// Construct the final SQL query with LIMIT and OFFSET
			paginatedSQL := rawSQL + " LIMIT " + limitPlaceholder + " OFFSET " + offsetPlaceholder
			// Append pagination parameters to existing args
			paginationArgs := []interface{}{pagination.GetLimit(), pagination.GetOffset()}
			finalArgs := append(args, paginationArgs...)
			// Apply raw SQL with pagination
			return db.Raw(paginatedSQL, finalArgs...)
		}
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}
