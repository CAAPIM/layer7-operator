### Retrieve Schema
```
$ npm install -g get-graphql-schema
$ NODE_TLS_REJECT_UNAUTHORIZED=0 get-graphql-schema -h "Authorization=Basic YWRtaW46N2xheWVy" https://localhost:9443/graphman > schema.graphql
```

### Generate Schema

```
$ go get github.com/Khan/genqlient@v0.7.0
$ go run github.com/Khan/genqlient@v0.7.0
```