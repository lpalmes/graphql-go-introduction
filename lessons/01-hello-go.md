# Hello GraphQL-Go!

If you know a bit about graphql you are in good shape to start, if you don't i will try to explain everything down the road.

## Here it is a quick rundown of the schema.

```
schema {
    query: Query
}

type Query  {
    hello: String!
}
```

Here the [schema](http://graphql.org/learn/schema/#the-query-and-mutation-types) defines the start of the graphql queries and mutations, we will only use Query for now.
This bit below is from [graphql.org](http://graphql.org)
>Every GraphQL service has a query type and may or may not have a mutation type. These types are the same as a regular object type, but they are special because they define the entry point of every GraphQL query.

The Query object here only contains one field, `hello: String!` which we will be able to use when working with our graphql server.

All of these Query fields should be resolved by us. So we use a `struct` to declare a main Resolver

This will be the main Resolver that we will use
```go
type Resolver struct {}
```

And this func will resolve our `hello: String!`.

```go
func (r *Resolver) Hello() string {
	return "world"
}
```
Notice that our func is exported by having the first letter beign uppercased, and our `func` returns a `string` because we declared `hello` to be a `String!`, The `!` declares that hello should not be null, so that's why we are able to return a `string` rather than a `*string`

That's everything that we will need! ᕕ( ᐛ )ᕗ

You can see in `main.go` that we also have
``` go
schemaFile, err := ioutil.ReadFile("schema.graphql")
schema, err = graphql.ParseSchema(string(schemaFile), &Resolver{})
```
to read our `schema.graphql` file and then convert it to a `graphql.Schema` type. We have to pass our `Resolver struct` so it can use our resolvers like `Hello()` in this case

Then we just pass it to a `http.Handle()` function and that's it
```go
http.Handle("/graphql", &relay.Handler{Schema: schema})
```

Keep in mind if you change the path you should change your path in graphiql.html script given that it points to your graphql handler

Now if you `go run main.go` you will be able to see that our server runs and if we go to `http://localhost:8080` you will be able to use graphiql

Just paste this into the left panel and click the play arrow at the top left (`cmd + enter` is faster for those who have to pickup their pet from the hair salon in 5 min)
```
{
    hello
}
```
This is the first query, and if you something like this at the right, all went pretty pretty good.
```
{
  "data": {
    "hello": "world"
  }
}
```