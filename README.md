# back-assignment-2-golang

Project based on GO

### Premise

I couldn't respect all the specification given.
The JSON-RPC method must be `{Type}.{Method}` with every JSON_RPC implementation found.

So you can't use `SearchNearestPharmacy` as method but you have to use `Search.NearestPharmacy`

### Requirements

- github.com/gorilla/mux
- github.com/gorilla/rpc
- Free 8081 port

### Installation

1) Clone the project
2) `cd back-assignment-2-golang`
3) Install dependencies
```
go get -u github.com/gorilla/mux
go get -u github.com/gorilla/rpc
```
4) Build it
`go build main.go`
   
### Run it

`./main`

### Use it

There is a Postman collection inside the root with the endpoints available with example requests in them. Play with it!