package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"io/ioutil"

	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
)

var PublicKey = []byte("secret")

type user struct {
	ID       string
	Name     string
	Mail     string
	Password string
}

type LoginInput struct {
	Mail     string
	Password string
}

var users = []user{
	{
		ID:       "1",
		Name:     "Example User",
		Mail:     "example@mail.com",
		Password: "example",
	},
	{
		ID:       "2",
		Name:     "Test User",
		Mail:     "test@mail.com",
		Password: "test",
	},
}

type viewerResolver struct {
	User user
}

// This will be use by our handler at /graphql
var schema *graphql.Schema

// This function runs at the start of the program
func init() {

	// We get the schema from the file, rather than having the schema inline here
	// I think will lead to better organizaiton of our own code
	schemaFile, err := ioutil.ReadFile("schema.graphql")
	if err != nil {
		// We will panic if we don't find the schema.graphql file in our server
		panic(err)
	}

	// We will use graphql-go library to parse our schema from "schema.graphql"
	// and the resolver is our struct that should fullfill everything in the Query
	// from our schema
	schema, err = graphql.ParseSchema(string(schemaFile), &Resolver{})
	if err != nil {
		panic(err)
	}
}

//Helper function
func checkToken(jwtToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return PublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return token, nil
	}

	return token, nil
}

func generateToken(user user) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(24)).Unix(), //This token will live for 24 hours
		"iat": time.Now().Unix(),
		"sub": user.ID,
	})
	tokenString, err := token.SignedString(PublicKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.Header.Get("Authorization")
		jwt, err := checkToken(token)
		if err != nil {
			fmt.Println(err)
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, "jwt", jwt)))
	})
}

func main() {
	// We will start a small server that reads our "graphiql.html" file and
	// responds with it, so we are able to have our own graphiql
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := ioutil.ReadFile("graphiql.html")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(page)
	}))

	// This is where our graphql server is handled, we declare "/graphql" as the route
	// where all our graphql requests will be directed to
	http.Handle("/graphql", auth(&relay.Handler{Schema: schema}))

	// We start the server by using ListenAndServe and we log if we have any error, hope not!
	fmt.Println("Listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Resolver struct is the main resolver  wich we will use to fullfill
// queries and mutations that our schema.graphql defines
type Resolver struct{}

// Hello function resolves to hello: String! in the Query object in our schema.graphql
func (r *Resolver) Hello() string {
	return "world"
}

func (r *Resolver) Login(args *struct {
	Input *LoginInput
}) (string, error) {
	for _, user := range users {
		if user.Mail == args.Input.Mail {
			if user.Password == args.Input.Password {
				token, err := generateToken(user)
				if err != nil {
					return "", err
				}
				return token, err
			} else {
				return "", errors.New("password is incorrect")
			}
		}
	}

	return "", errors.New("User not found")
}

func (r *Resolver) Viewer(ctx context.Context, args *struct {
	Token *string
}) (*viewerResolver, error) {

	token := ctx.Value("jwt").(*jwt.Token)
	if token == nil && args.Token == nil {
		return nil, errors.New("There needs to be a token in the Authorization header or viewer input")
	}

	if token == nil && args.Token != nil {
		viewerToken, err := checkToken(*args.Token)
		if err != nil {
			return nil, err
		}

		token = viewerToken
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	id := claims["sub"].(string)

	var user user

	for _, u := range users {
		if id == string(u.ID) {
			user = u
		}
	}

	return &viewerResolver{
		User: user,
	}, nil
}

func (v *viewerResolver) ID() graphql.ID {
	return graphql.ID(v.User.ID)
}

func (v *viewerResolver) Name() string {
	return v.User.Name
}

func (v *viewerResolver) Mail() string {
	return v.User.Mail
}
