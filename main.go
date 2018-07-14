package main

//ToDo
//Import Library
//Create Struct for Menu
//Initiation Variable
//Initiation Menu Type for GraphQL
//Create Root Mutation
//Create Root Query
// Define GraphQL Schema
// Create Main Function


//Import Library
import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
)

//Create Struct for Menu
type Menu struct {
	Name        string `json:"name"`
	Main        string `json:"main"`
	Method      string `json:"method"`
	Additional  string `json:"additional"`
}

//Initiation Variable
var MenuList []Menu

func init() {
	menu1 := Menu{Name: "nasi-goreng-ayam", Main: "Nasi", Method: "Goreng", Additional: "Ayam"}
	menu2 := Menu{Name: "nasi-goreng-teri", Main: "Nasi", Method: "Goreng", Additional: "Teri"}
	menu3 := Menu{Name: "nasi-goreng-sosis", Main: "Nasi", Method: "Goreng", Additional: "Sosis"}

	menu4 := Menu{Name: "mie-ayam-bakso", Main: "Mie", Method: "Ayam", Additional: "Bakso"}
	menu5 := Menu{Name: "mie-ayam-ceker", Main: "Mie", Method: "Ayam", Additional: "Ceker"}
	menu6 := Menu{Name: "mie-ayam-pangsit", Main: "Mie", Method: "Ayam", Additional: "Pangsit"}

	MenuList = append(MenuList, menu1, menu2, menu3, menu4, menu5, menu6)
}

//Initiation Menu Type for GraphQL
var menuType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Menu",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"main": &graphql.Field{
			Type: graphql.String,
		},
		"method": &graphql.Field{
			Type: graphql.String,
		},
		"additional": &graphql.Field{
			Type: graphql.String,
		},
	},
})

//Create Root Mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		/*
			curl -g 'http://localhost:8080/graphql?query=mutation+_{createTodo(text:"My+new+todo"){name,main,method,additional}}'
		*/
		"createMenu": &graphql.Field{
			Type:        menuType, // the return type for this field
			Description: "Create new menu",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"main": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"method": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"additional": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				// marshall and cast the argument value
				name, _ := params.Args["name"].(string)
				main, _ := params.Args["main"].(string)
				method, _ := params.Args["method"].(string)
				additional, _ := params.Args["additional"].(string)

				// perform mutation operation here
				// for e.g. create a Menu and save to DB.
				newMenu := Menu{
					Name: name,
					Main: main,
					Method: method,
					Additional: additional,
				}

				MenuList = append(MenuList, newMenu)

				return newMenu, nil
			},
		},
		/*
			curl -g 'http://localhost:8080/graphql?query=mutation+_{updateMenu(id:"a",done:true){id,text,done}}'
		*/
		"updateMenu": &graphql.Field{
			Type:        menuType, // the return type for this field
			Description: "Update existing menu",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"main": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"method": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"additional": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// marshall and cast the argument value
				name, _ := params.Args["name"].(string)
				main, _ := params.Args["main"].(string)
				method, _ := params.Args["method"].(string)
				additional, _ := params.Args["additional"].(string)
				affectedMenu := Menu{}

				// Search list for todo with id and change the done variable
				for i := 0; i < len(MenuList); i++ {
					if MenuList[i].Name == name {
						MenuList[i].Main = main
						MenuList[i].Method = method
						MenuList[i].Additional = additional
						// Assign updated menu so we can return it
						affectedMenu = MenuList[i]
						break
					}
				}
				// Return affected todo
				return affectedMenu, nil
			},
		},
	},
})

//Create Root Query
// we just define a trivial example here, since root query is required.
// Test with curl
// curl -g 'http://localhost:8080/graphql?query={lastTodo{id,text,done}}'
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{

		//curl -g 'http://localhost:8080/graphql?query={menu(name:"nasi-goreng-ayam"){name,main,method,additional}}'

		"menu": &graphql.Field{
			Type:        menuType,
			Description: "Get single menu",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"main": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"method": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"additional": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				name, isOK := params.Args["name"].(string)
				if isOK {
					// Search for el with name
					for _, menu := range MenuList {
						if menu.Name == name {
							return menu, nil
						}
					}
				}

				return Menu{}, nil
			},
		},

		"lastMenu": &graphql.Field{
			Type:        menuType,
			Description: "Last menu added",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return MenuList[len(MenuList)-1], nil
			},
		},

		//curl -g 'http://localhost:8080/graphql?query={menuList{name,main,method,additional}}'

		"menuList": &graphql.Field{
			Type:        graphql.NewList(menuType),
			Description: "List of menus",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return MenuList, nil
			},
		},
	},
})

// Define GraphQL Schema
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

//Create Main Function
func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	// Display some basic instructions
	fmt.Println("Now server is running on port 8080")
	fmt.Println("Get single menu: curl -g 'http://localhost:8080/graphql?query={menu(name:\"nasi-goreng-ayam\"){name,main,method,additional}}'")
	fmt.Println("Create new menu: curl -g 'http://localhost:8080/graphql?query=mutation+_{createMenu(name:\"nasi-goreng-udang\",main:\"Nasi\",method:\"Goreng\",additional:\"Udang\"){name,main,method,additional}}'")
	fmt.Println("Update menu: curl -g 'http://localhost:8080/graphql?query=mutation+_{createMenu(name:\"nasi-goreng-udang\",main:\"Nasi\",method:\"Goreng\",additional:\"Udang\"){name,main,method,additional}}'")
	fmt.Println("Load menu list: curl -g 'http://localhost:8080/graphql?query={menuList{name,main,method,additional}}'")
	fmt.Println("Access the web app via browser at 'http://localhost:8080'")

	http.ListenAndServe(":8080", nil)
}