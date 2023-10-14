### Generate Schema

```
$ go get github.com/Khan/genqlient@v0.5.0
$ go run github.com/Khan/genqlient@v0.5.0
```

### Retrieve Schema
```
$ npm install -g get-graphql-schema
$ NODE_TLS_REJECT_UNAUTHORIZED=0 get-graphql-schema -h "Authorization=Basic YWRtaW46N2xheWVy" https://localhost:9443/graphman > schema.graphql
```

### NOTE
generated.go treats EncassArgInput.GuiPrompt as a string regardless of the schema... has to be changed manually. 