# k8s tdd in golang with kind, pulumi & k8s.io/client-go 

## Run
```
kind create cluster
pulumi stack init k8s-go-tdd
pulumi up --yes
go test ./...
```

## Tear down
```
pulumi destroy
pulumi stack rm k8s-go-tdd
kind delete cluster
```