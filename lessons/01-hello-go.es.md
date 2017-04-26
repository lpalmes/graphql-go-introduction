# Hola GraphQL-Go!


Si sabes un poco sobre graphql podes empezar tranquilamente, sino voy a intentar explicarlo mientras avanzemos.

## Este es nuestro schema

```
schema {
    query: Query
}

type Query  {
    hello: String!
}
```

Aca el [schema](http://graphql.org/learn/schema/#the-query-and-mutation-types) define el comienzo de los queries y las mutations, por ahora solamente vamos a usar Query.
Voy a intentar traducir este texto de [graphql.org](http://graphql.org).
>Cada servicio de GraphQL tiene un objeto tipo Query y puede tener o no un objeto tipo Mutation. Estos tipos de objetos son lo mismo que objetos regulares, pero en este caso son especiales por que definen la entrada a cada query de GraphQL. 

El objeto Query aca tiene una sola propiedad, `hello: String!` la cual vamos a utilizar cuando usemos nuestro servidor graphql.

Todos las propiedades de Query deben ser resueltas por nosotros. Por eso vamos a usar un `struct` para declarar un `Resolver`

Este va a ser el struct principal que vamos a usar para resolver todos los Query y Mutations que tengamos.
```go
type Resolver struct {}
```

Y esta funcion va a resolver nuestro `hello: String!`.

```go
func (r *Resolver) Hello() string {
	return "world"
}
```
Fijate que nuestra funcion esta exportada al tener la primera letra en mayuscula, y ademas devuelve un `string` por que declaramos `hello` de tipo `String!`, el `!` declara que hello no debe ser null, por eso podemos devolver un `string` en vez de un `*string`

Esto es todo lo que vamos a necesitar! ᕕ( ᐛ )ᕗ

Se puede ver que en `main.go` tenemos
``` go
schemaFile, err := ioutil.ReadFile("schema.graphql")
schema, err = graphql.ParseSchema(string(schemaFile), &Resolver{})
```
para leer nuestro archivo `schema.graphql` que va a convertirlo en un struct tipo `graphqlSchema`. Le tenemos que pasar nuestro `Resolver struct` para que resuelva los queries como `Hello()` en este caso.

Despues solamente se lo pasamos a nuestra funcion `http.Handle()` y eso es todo
```go
http.Handle("/graphql", &relay.Handler{Schema: schema})
```

Tene en cuenta que si cambias la direccion de nuestro `Handle` tambien tenes que cambiarla en el archivo `graphiql.html` para que apunte al servidor apropiadamente


Ahora si corres `go run main.go` vas a poder ver que nuestro servidor esta funcionando en `http://localhost:8080` y vamos a poder usar graphiql

Solamente copia esto en el panel izquierdo y dale click a la flecha de arriba (mac: `cmd + enter` windows: `ctrl + enter` si estas muy apurado)
```
{
    hello
}
```
Este es nuestro primer query, y si a la derecha hay algo como esto, parece que todo anduvo perfectamente.
```
{
  "data": {
    "hello": "world"
  }
}
```