package processor

import (
	"errors"
	"fmt"
	"go-play-graphql/entity"

	"github.com/graphql-go/graphql"
)

type RedirectProcessor struct {
	Storage *entity.DB
}

func (rp *RedirectProcessor) ExecuteQuery(query string) *graphql.Result {
	schema, _ := rp.schema()
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func (rp *RedirectProcessor) schema() (graphql.Schema, error) {
	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    rp.queryType(),
			Mutation: rp.mutationType(),
		},
	)
}

func redirectType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Redirect",
			Fields: graphql.Fields{
				"from": &graphql.Field{
					Type: graphql.String,
				},
				"to": &graphql.Field{
					Type: graphql.String,
				},
			},
		})
}

func (rp *RedirectProcessor) queryType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"redirect": &graphql.Field{
					Type: redirectType(),
					Args: graphql.FieldConfigArgument{
						"from": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: rp.get,
				},
				"redirects": &graphql.Field{
					Type: graphql.NewList(redirectType()),
					Args: graphql.FieldConfigArgument{
						"offset": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"limit": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: rp.getAll,
				},
			},
		})
}

func (rp *RedirectProcessor) mutationType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createRedirect": &graphql.Field{
				Type: redirectType(),
				Args: graphql.FieldConfigArgument{
					"from": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"to": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: rp.createRedirect,
			},
			"deleteRedirect": &graphql.Field{
				Type: redirectType(),
				Args: graphql.FieldConfigArgument{
					"from": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: rp.deleteRedirect,
			},
		},
	})
}

func (rp *RedirectProcessor) get(p graphql.ResolveParams) (interface{}, error) {
	from, ok := p.Args["from"].(string)
	if !ok {
		return nil, errors.New("Unable to parse from")
	}
	res, err := rp.Storage.Get(from)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (rp *RedirectProcessor) getAll(p graphql.ResolveParams) (interface{}, error) {
	offset, ok := p.Args["offset"].(int)
	if !ok {
		return nil, errors.New("Unable to parse offset")
	}
	limit, ok := p.Args["limit"].(int)
	if !ok {
		return nil, errors.New("Unable to parse limit")
	}
	result, err := rp.Storage.GetAll(offset, limit)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (rp *RedirectProcessor) createRedirect(params graphql.ResolveParams) (interface{}, error) {
	from, _ := params.Args["from"].(string)
	to, _ := params.Args["to"].(string)
	return rp.Storage.Save(&entity.Redirect{From: from, To: to})
}

func (rp *RedirectProcessor) deleteRedirect(params graphql.ResolveParams) (interface{}, error) {
	from, _ := params.Args["from"].(string)
	err := rp.Storage.Delete(from)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
