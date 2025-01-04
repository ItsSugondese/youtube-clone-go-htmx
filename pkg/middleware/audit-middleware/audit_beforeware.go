package audit_middleware

import (
	"youtube-clone/global/marker"
	"github.com/gin-gonic/gin"
	user_data "youtube-clone/pkg/utils/user-data"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"time"
)

func RegisterCallbacks(db *gorm.DB) error {
	db.Callback().Create().Before("gorm:create").Register("custom_plugin:create_audit_log", createAuditLog)
	db.Callback().Update().Before("gorm:update").Register("custom_plugin:update_audit_log", updateAuditLog)
	//db.Callback().Delete().Before("gorm:delete").Register("custom_plugin:delete_audit_log", deleteAuditLog)
	return nil
}

func createAuditLog(db *gorm.DB) {
	if model, ok := db.Statement.Model.(marker.Auditable); ok {
		if !model.HasAuditModel() {
			return
		}
	} else {
		return
	}

	if db.Statement.Schema != nil && db.Statement.Schema.Table == "audit_logs" || db.Error != nil ||
		db.Statement.Schema.Table == "user_role" {
		return
	}

	userId := ""
	if ginCtx, ok := db.Statement.Context.(*gin.Context); ok {
		// Now you have the gin.Context and can use it as needed
		userId, _ = user_data.GetUserIdContext(ginCtx) // Example: retrieve a value from gin.Context
	}
	var ifNil *string

	if userId == "" {
		ifNil = nil
	} else {
		ifNil = &userId
	}

	updateAuditModelFields(db.Statement.Model, ifNil, "POST")
}

func updateAuditLog(db *gorm.DB) {
	if model, ok := db.Statement.Model.(marker.Auditable); ok {
		if !model.HasAuditModel() {
			return
		}
	} else {
		return
	}

	userId := ""
	if ginCtx, ok := db.Statement.Context.(*gin.Context); ok {
		// Now you have the gin.Context and can use it as needed
		userId, _ = user_data.GetUserIdContext(ginCtx) // Example: retrieve a value from gin.Context
	}

	var ifNil *string

	if userId == "" {
		ifNil = nil
	} else {
		ifNil = &userId
	}

	model := db.Statement.Model
	value := reflect.ValueOf(model).Elem()

	field := value.FieldByName("IsDeleting")

	boolValue := field.Bool()

	if boolValue {
		deleteAuditLog(db)
		return
	}
	updateAuditModelFields(model, ifNil, "PUT")

}

func deleteAuditLog(db *gorm.DB) {
	if model, ok := db.Statement.Model.(marker.Auditable); ok {
		if !model.HasAuditModel() {
			return
		}
	} else {
		return
	}

	if db.Statement.Schema != nil && db.Statement.Schema.Table == "audit_logs" || db.Error != nil {
		return
	}

	userId := ""
	if ginCtx, ok := db.Statement.Context.(*gin.Context); ok {
		// Now you have the gin.Context and can use it as needed
		userId, _ = user_data.GetUserIdContext(ginCtx) // Example: retrieve a value from gin.Context
	}
	var ifNil *string

	if userId == "" {
		ifNil = nil
	} else {
		ifNil = &userId
	}

	updateAuditModelFields(db.Statement.Model, ifNil, "DELETE")
}

func updateAuditModelFields(model interface{}, userId *string, saveType string) {
	value := reflect.ValueOf(model)

	// If model is a pointer, get the underlying value
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Check if the value is a struct
	if value.Kind() != reflect.Struct {
		fmt.Println("Expected a struct, got:", value.Kind())
		return
	}

	if saveType == "POST" {
		// Set CreatedBy field
		createdByField := value.FieldByName("CreatedBy")
		if createdByField.IsValid() && createdByField.CanSet() {
			createdByField.Set(reflect.ValueOf(userId))

		}

		// Set CreatedAt field
		createdAtField := value.FieldByName("CreatedAt")
		if createdAtField.IsValid() && createdAtField.CanSet() {
			createdAtField.Set(reflect.ValueOf(time.Now()))
		}
	}

	if saveType == "PUT" {

		// Set UpdatedAt field
		updatedAtField := value.FieldByName("UpdatedAt")
		if updatedAtField.IsValid() && updatedAtField.CanSet() {
			updatedAtField.Set(reflect.ValueOf(time.Now()))
		}

		// Set UpdatedBy field
		updatedByField := value.FieldByName("UpdatedBy")
		if updatedByField.IsValid() && updatedByField.CanSet() {

			updatedByField.Set(reflect.ValueOf(userId))
		}
	}

	if saveType == "DELETE" {
		// Set DeletedBy field
		deletedByField := value.FieldByName("DeletedBy")
		if deletedByField.IsValid() && deletedByField.CanSet() {
			deletedByField.Set(reflect.ValueOf(userId))
		}

		deletedAtField := value.FieldByName("DeletedAt")
		if deletedAtField.IsValid() && deletedAtField.CanSet() {
			// Create an instance of gorm.DeletedAt
			deletedAtValue := gorm.DeletedAt{
				Time:  time.Now(),
				Valid: true,
			}
			// Set the field value
			deletedAtField.Set(reflect.ValueOf(deletedAtValue))
		}
	}
}
