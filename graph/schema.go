package graph

import (
	"fmt" // Add this import
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
	"scrumer-backend/models" // Add this import
)

type Resolver struct {
	DB *gorm.DB
}

func NewSchema(db *gorm.DB) (graphql.Schema, error) {
	resolver := &Resolver{DB: db} // Keep this line for future use

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: graphql.NewObject(
				graphql.ObjectConfig{
					Name: "Query",
					Fields: graphql.Fields{
						"hello": &graphql.Field{
							Type: graphql.String,
							Resolve: func(p graphql.ResolveParams) (interface{}, error) {
								return "world", nil
							},
						},
						"users": &graphql.Field{
							Type: graphql.NewList(userType),
							Resolve: resolver.GetUsers,
						},
					},
				},
			),
		},
	)
	if err != nil {
		return graphql.Schema{}, err
	}
	return schema, nil
}

// GetUsers resolver to fetch all users
func (r *Resolver) GetUsers(p graphql.ResolveParams) (interface{}, error) {
	var users []models.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}
