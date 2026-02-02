package graph

import (
	"fmt"
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"scrumer-backend/models"
)

type Resolver struct {
	DB GormDB
}

func NewSchema(db GormDB) (graphql.Schema, error) {
	resolver := &Resolver{DB: db}

	// Define dateTimeType locally
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

	// Declare userType and projectType as vars here to handle circular dependency
	var userType *graphql.Object
	var projectType *graphql.Object
	var sprintType *graphql.Object
	var taskType *graphql.Object

	// Define userType locally
	userType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "User",
			Fields: graphql.FieldsThunk(func() graphql.Fields {
				return graphql.Fields{
					"id": &graphql.Field{
						Type: graphql.ID,
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							u := p.Source.(models.User)
							return strconv.FormatUint(uint64(u.ID), 10), nil
						},
					},
					"username":  &graphql.Field{Type: graphql.String},
					"email":     &graphql.Field{Type: graphql.String},
					"createdAt": &graphql.Field{Type: dateTimeType},
					"updatedAt": &graphql.Field{Type: dateTimeType},
					"projects":  &graphql.Field{
						Type: graphql.NewList(projectType),
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							user, ok := p.Source.(models.User)
							if !ok {
								return nil, fmt.Errorf("invalid user source type")
							}
							var projects []models.Project
							if err := resolver.DB.Model(&user).Association("Projects").Find(&projects); err != nil {
								return nil, fmt.Errorf("failed to get projects for user: %w", err)
							}
							return projects, nil
						},
					},
				}
			}),
		},
	)

	// Define projectType locally
	projectType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Project",
			Fields: graphql.FieldsThunk(func() graphql.Fields {
				return graphql.Fields{
					"id": &graphql.Field{
						Type: graphql.ID,
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							proj := p.Source.(models.Project)
							return strconv.FormatUint(uint64(proj.ID), 10), nil
						},
					},
					"name":        &graphql.Field{Type: graphql.String},
					"description": &graphql.Field{Type: graphql.String},
					"createdAt":   &graphql.Field{Type: dateTimeType},
					"updatedAt":   &graphql.Field{Type: dateTimeType},
					"users":       &graphql.Field{
						Type: graphql.NewList(userType),
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							project, ok := p.Source.(models.Project)
							if !ok {
								return nil, fmt.Errorf("invalid project source type")
							}
							var users []models.User
							if err := resolver.DB.Model(&project).Association("Users").Find(&users); err != nil {
								return nil, fmt.Errorf("failed to get users for project: %w", err)
							}
							return users, nil
						},
					},
				}
			}),
		},
	)

	// Define sprintType locally
	sprintType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Sprint",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.ID,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						s := p.Source.(models.Sprint)
						return strconv.FormatUint(uint64(s.ID), 10), nil
					},
				},
				"name":      &graphql.Field{Type: graphql.String},
				"startDate": &graphql.Field{Type: dateTimeType},
				"endDate":   &graphql.Field{Type: dateTimeType},
				"createdAt": &graphql.Field{Type: dateTimeType},
				"updatedAt": &graphql.Field{Type: dateTimeType},
			},
		},
	)

	// Define taskType locally
	taskType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Task",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.ID,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						t := p.Source.(models.Task)
						return strconv.FormatUint(uint64(t.ID), 10), nil
					},
				},
				"title":       &graphql.Field{Type: graphql.String},
				"description": &graphql.Field{Type: graphql.String},
				"status":      &graphql.Field{Type: graphql.String},
				"priority":    &graphql.Field{Type: graphql.String},
				"createdAt": &graphql.Field{Type: dateTimeType},
				"updatedAt": &graphql.Field{Type: dateTimeType},
			},
		},
	)

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
							Type:    graphql.NewList(userType),
							Resolve: resolver.GetUsers,
						},
						"project": &graphql.Field{
							Type:        projectType,
							Description: "Get single project",
							Args: graphql.FieldConfigArgument{
								"id": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.ID),
								},
							},
							Resolve: resolver.GetProject,
						},
						"projects": &graphql.Field{
							Type:        graphql.NewList(projectType),
							Description: "Get all projects",
							Resolve:     resolver.GetProjects,
						},
						"sprint": &graphql.Field{
							Type:        sprintType,
							Description: "Get single sprint",
							Args: graphql.FieldConfigArgument{
								"id": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.ID),
								},
							},
							Resolve: resolver.GetSprint,
						},
						"sprints": &graphql.Field{
							Type:    graphql.NewList(sprintType),
							Description: "Get all sprints",
							Resolve: resolver.GetSprints,
						},
						"task": &graphql.Field{
							Type:        taskType,
							Description: "Get single task",
							Args: graphql.FieldConfigArgument{
								"id": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.ID),
								},
							},
							Resolve: resolver.GetTask,
						},
						"tasks": &graphql.Field{
							Type:    graphql.NewList(taskType),
							Description: "Get all tasks",
							Resolve: resolver.GetTasks,
						},
					},
				},
			),
			Mutation: graphql.NewObject(
				graphql.ObjectConfig{
					Name: "Mutation",
					Fields: graphql.Fields{
						"createUser": &graphql.Field{
							Type: userType,
							Args: graphql.FieldConfigArgument{
								"username": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.String),
								},
								"email": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.String),
								},
								"password": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.String),
								},
							},
							Resolve: resolver.CreateUser,
						},
						"updateUser": &graphql.Field{
							Type: userType,
							Args: graphql.FieldConfigArgument{
								"id": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.ID),
								},
								"username": &graphql.ArgumentConfig{
									Type: graphql.String,
								},
								"email": &graphql.ArgumentConfig{
									Type: graphql.String,
								},
								"password": &graphql.ArgumentConfig{
									Type: graphql.String,
								},
							},
							Resolve: resolver.UpdateUser,
						},
						"deleteUser": &graphql.Field{
							Type: graphql.Boolean,
							Args: graphql.FieldConfigArgument{
								"id": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.ID),
								},
							},
							Resolve: resolver.DeleteUser,
						},
						"createProject": &graphql.Field{
							Type: projectType,
							Args: graphql.FieldConfigArgument{
								"name": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.String),
								},
								"description": &graphql.ArgumentConfig{
									Type: graphql.String,
								},
							},
							Resolve: resolver.CreateProject,
						},
						"updateProject": &graphql.Field{
							Type: projectType,
							Args: graphql.FieldConfigArgument{
								"id": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.ID),
								},
								"name": &graphql.ArgumentConfig{
									Type: graphql.String,
								},
								"description": &graphql.ArgumentConfig{
									Type: graphql.String,
								},
							},
							Resolve: resolver.UpdateProject,
						},
						"deleteProject": &graphql.Field{
							Type: graphql.Boolean,
							Args: graphql.FieldConfigArgument{
								"id": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.ID),
								},
							},
							Resolve: resolver.DeleteProject,
						},
						"addUserToProject": &graphql.Field{
							Type: projectType,
							Args: graphql.FieldConfigArgument{
								"projectId": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.ID),
								},
								"userId": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.ID),
								},
							},
							Resolve: resolver.AddUserToProject,
						},
						"removeUserFromProject": &graphql.Field{
							Type: projectType,
							Args: graphql.FieldConfigArgument{
								"projectId": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.ID),
								},
								"userId": &graphql.ArgumentConfig{
									Type: graphql.NewNonNull(graphql.ID),
								},
							},
							Resolve: resolver.RemoveUserFromProject,
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

// GetUsers resolver to fetch all users (This will remain as a method on Resolver)
func (r *Resolver) GetUsers(p graphql.ResolveParams) (interface{}, error) {
	var users []models.User
	// Preload Projects for users to avoid N+1 problem
	if err := r.DB.Preload("Projects").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

// GetProject resolver to fetch a single project by ID
func (r *Resolver) GetProject(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid project ID")
	}
	var project models.Project
	if err := r.DB.Preload("Users").First(&project, id).Error; err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}
	return project, nil
}

// GetProjects resolver to fetch all projects
func (r *Resolver) GetProjects(p graphql.ResolveParams) (interface{}, error) {
	var projects []models.Project
	if err := r.DB.Preload("Users").Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	return projects, nil
}

// GetSprint resolver to fetch a single sprint by ID
func (r *Resolver) GetSprint(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid sprint ID")
	}
	var sprint models.Sprint
	if err := r.DB.First(&sprint, id).Error; err != nil {
		return nil, fmt.Errorf("sprint not found: %w", err)
	}
	return sprint, nil
}

// GetSprints resolver to fetch all sprints
func (r *Resolver) GetSprints(p graphql.ResolveParams) (interface{}, error) {
	var sprints []models.Sprint
	if err := r.DB.Find(&sprints).Error; err != nil {
		return nil, fmt.Errorf("failed to get sprints: %w", err)
	}
	return sprints, nil
}

// GetTask resolver to fetch a single task by ID
func (r *Resolver) GetTask(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid task ID")
	}
	var task models.Task
	if err := r.DB.First(&task, id).Error; err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}
	return task, nil
}

// GetTasks resolver to fetch all tasks
func (r *Resolver) GetTasks(p graphql.ResolveParams) (interface{}, error) {
	var tasks []models.Task
	if err := r.DB.Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	return tasks, nil
}

// CreateUser resolver
func (r *Resolver) CreateUser(p graphql.ResolveParams) (interface{}, error) {
	username, _ := p.Args["username"].(string)
	email, _ := p.Args["email"].(string)
	password, _ := p.Args["password"].(string)

	user := models.User{
		Username: username,
		Email:    email,
		Password: password, // In a real application, hash this password!
	}
	if err := r.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// UpdateUser resolver
func (r *Resolver) UpdateUser(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID")
	}

	var user models.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if username, ok := p.Args["username"].(string); ok {
		user.Username = username
	}
	if email, ok := p.Args["email"].(string); ok {
		user.Email = email
	}
	if password, ok := p.Args["password"].(string); ok {
		user.Password = password // In a real application, hash this password!
	}

	if err := r.DB.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

// DeleteUser resolver
func (r *Resolver) DeleteUser(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID")
	}

	if err := r.DB.Delete(&models.User{}, id).Error; err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}
	return true, nil // Return true for success
}

// CreateProject resolver
func (r *Resolver) CreateProject(p graphql.ResolveParams) (interface{}, error) {
	name, _ := p.Args["name"].(string)
	description, _ := p.Args["description"].(string)

	project := models.Project{
		Name:        name,
		Description: description,
	}
	if err := r.DB.Create(&project).Error; err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	return project, nil
}

// UpdateProject resolver
func (r *Resolver) UpdateProject(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid project ID")
	}

	var project models.Project
	if err := r.DB.First(&project, id).Error; err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	if name, ok := p.Args["name"].(string); ok {
		project.Name = name
	}
	if description, ok := p.Args["description"].(string); ok {
		project.Description = description
	}

	if err := r.DB.Save(&project).Error; err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}
	return project, nil
}

// DeleteProject resolver
func (r *Resolver) DeleteProject(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid project ID")
	}

	if err := r.DB.Delete(&models.Project{}, id).Error; err != nil {
		return nil, fmt.Errorf("failed to delete project: %w", err)
	}
	return true, nil // Return true for success
}

// AddUserToProject resolver
func (r *Resolver) AddUserToProject(p graphql.ResolveParams) (interface{}, error) {
	projectID, ok := p.Args["projectId"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid project ID")
	}
	userID, ok := p.Args["userId"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID")
	}

	var project models.Project
	if err := r.DB.First(&project, projectID).Error; err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}
	var user models.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if err := r.DB.Model(&project).Association("Users").Append(&user); err != nil {
		return nil, fmt.Errorf("failed to add user to project: %w", err)
	}
	return project, nil
}

// RemoveUserFromProject resolver
func (r *Resolver) RemoveUserFromProject(p graphql.ResolveParams) (interface{}, error) {
	projectID, ok := p.Args["projectId"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid project ID")
	}
	userID, ok := p.Args["userId"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID")
	}

	var project models.Project
	if err := r.DB.First(&project, projectID).Error; err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}
	var user models.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if err := r.DB.Model(&project).Association("Users").Delete(&user); err != nil {
		return nil, fmt.Errorf("failed to remove user from project: %w", err)
	}
	return project, nil
}
