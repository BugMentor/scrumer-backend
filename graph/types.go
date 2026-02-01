package graph

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

var dateTimeType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "DateTime",
	Description: "DateTime custom scalar type",
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case time.Time:
			return value.Format(time.RFC3339)
		case *time.Time:
			if value == nil {
				return nil
			}
			return value.Format(time.RFC3339)
		default:
			return nil
		}
	},
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case string:
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return nil
			}
			return t
		default:
			return nil
		}
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			t, err := time.Parse(time.RFC3339, valueAST.Value)
			if err != nil {
				return nil
			}
			return t
		default:
			return nil
		}
	},
})

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":       &graphql.Field{Type: graphql.Int},
			"username": &graphql.Field{Type: graphql.String},
			"email":    &graphql.Field{Type: graphql.String},
			"createdAt": &graphql.Field{Type: dateTimeType},
			"updatedAt": &graphql.Field{Type: dateTimeType},
		},
	},
)

var projectType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Project",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.Int},
			"name":        &graphql.Field{Type: graphql.String},
			"description": &graphql.Field{Type: graphql.String},
			"createdAt": &graphql.Field{Type: dateTimeType},
			"updatedAt": &graphql.Field{Type: dateTimeType},
		},
	},
)

var sprintType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Sprint",
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.Int},
			"name":      &graphql.Field{Type: graphql.String},
			"startDate": &graphql.Field{Type: dateTimeType},
			"endDate":   &graphql.Field{Type: dateTimeType},
			"createdAt": &graphql.Field{Type: dateTimeType},
			"updatedAt": &graphql.Field{Type: dateTimeType},
		},
	},
)

var taskType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Task",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.Int},
			"title":       &graphql.Field{Type: graphql.String},
			"description": &graphql.Field{Type: graphql.String},
			"status":      &graphql.Field{Type: graphql.String},
			"priority":    &graphql.Field{Type: graphql.String},
			"createdAt": &graphql.Field{Type: dateTimeType},
			"updatedAt": &graphql.Field{Type: dateTimeType},
		},
	},
)