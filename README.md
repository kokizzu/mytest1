# TiDB test

## dependencies

```bash
go mod tidy
docker-compose up -d
sudo chown -R `whoami` */ 
```

## manual migration (run automatically on init)

```bash
go get -u -v github.com/naoina/migu/cmd/migu
 
migu sync -t mysql -u root -h 127.0.0.1 -P 4000 test schema.go
```
