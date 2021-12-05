# TiDB example using golang, sqlx, and migu

this example using `IN` AND `LIKE` query.

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

## result

```bash
go run .

Time for 10000 operations:

                  Insert Max    145913133 ns =    145913 µs =    145.913 ms =  0.1459 s
                  Insert Min    145913133 ns =    145913 µs =    145.913 ms =  0.1459 s
                Insert Total    145913133 ns =    145913 µs =    145.913 ms =  0.1459 s
                  Insert Avg        14591 ns =        14 µs =      0.015 ms =  0.0000 s

               []*Struct Max     16504507 ns =     16504 µs =     16.505 ms =  0.0165 s
               []*Struct Min       705456 ns =       705 µs =      0.705 ms =  0.0007 s
             []*Struct Total  10867411840 ns =  10867411 µs =  10867.412 ms = 10.8674 s
               []*Struct Avg      1086741 ns =      1086 µs =      1.087 ms =  0.0011 s

                   []Map Max      6682211 ns =      6682 µs =      6.682 ms =  0.0067 s
                   []Map Min       714383 ns =       714 µs =      0.714 ms =  0.0007 s
                 []Map Total  10683510979 ns =  10683510 µs =  10683.511 ms = 10.6835 s
                   []Map Avg      1068351 ns =      1068 µs =      1.068 ms =  0.0011 s

             []MapManual Max      4031078 ns =      4031 µs =      4.031 ms =  0.0040 s
             []MapManual Min       674298 ns =       674 µs =      0.674 ms =  0.0007 s
           []MapManual Total  10560956893 ns =  10560956 µs =  10560.957 ms = 10.5610 s
             []MapManual Avg      1056095 ns =      1056 µs =      1.056 ms =  0.0011 s

                 []Slice Max     12833500 ns =     12833 µs =     12.834 ms =  0.0128 s
                 []Slice Min       688754 ns =       688 µs =      0.689 ms =  0.0007 s
               []Slice Total  10560311375 ns =  10560311 µs =  10560.311 ms = 10.5603 s
                 []Slice Avg      1056031 ns =      1056 µs =      1.056 ms =  0.0011 s

           []SliceManual Max      9376333 ns =      9376 µs =      9.376 ms =  0.0094 s
           []SliceManual Min       684256 ns =       684 µs =      0.684 ms =  0.0007 s
         []SliceManual Total  10724436604 ns =  10724436 µs =  10724.437 ms = 10.7244 s
           []SliceManual Avg      1072443 ns =      1072 µs =      1.072 ms =  0.0011 s

                []Struct Max     18099372 ns =     18099 µs =     18.099 ms =  0.0181 s
                []Struct Min       682513 ns =       682 µs =      0.683 ms =  0.0007 s
              []Struct Total  10915536051 ns =  10915536 µs =  10915.536 ms = 10.9155 s
                []Struct Avg      1091553 ns =      1091 µs =      1.092 ms =  0.0011 s

Done in  1m4.537577702s
```
