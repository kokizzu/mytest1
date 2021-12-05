# TiDB example using golang, sqlx, and migu

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
Time for 10000 operations:

               []*Struct Max     14539728 ns =     14539 µs =     14.54 ms =  0.0145 s
               []*Struct Min       712454 ns =       712 µs =      0.71 ms =  0.0007 s
             []*Struct Total  10928911606 ns =  10928911 µs =  10928.91 ms = 10.9289 s
               []*Struct Avg      1092891 ns =      1092 µs =      1.09 ms =  0.0011 s

                   []Map Max      7049101 ns =      7049 µs =      7.05 ms =  0.0070 s
                   []Map Min       699280 ns =       699 µs =      0.70 ms =  0.0007 s
                 []Map Total  11110398351 ns =  11110398 µs =  11110.40 ms = 11.1104 s
                   []Map Avg      1111039 ns =      1111 µs =      1.11 ms =  0.0011 s

             []MapManual Max      7420427 ns =      7420 µs =      7.42 ms =  0.0074 s
             []MapManual Min       731329 ns =       731 µs =      0.73 ms =  0.0007 s
           []MapManual Total  11018502810 ns =  11018502 µs =  11018.50 ms = 11.0185 s
             []MapManual Avg      1101850 ns =      1101 µs =      1.10 ms =  0.0011 s

                 []Slice Max     15544140 ns =     15544 µs =     15.54 ms =  0.0155 s
                 []Slice Min       706433 ns =       706 µs =      0.71 ms =  0.0007 s
               []Slice Total  10990063585 ns =  10990063 µs =  10990.06 ms = 10.9901 s
                 []Slice Avg      1099006 ns =      1099 µs =      1.10 ms =  0.0011 s

           []SliceManual Max     10644994 ns =     10644 µs =     10.64 ms =  0.0106 s
           []SliceManual Min       717323 ns =       717 µs =      0.72 ms =  0.0007 s
         []SliceManual Total  10737251440 ns =  10737251 µs =  10737.25 ms = 10.7373 s
           []SliceManual Avg      1073725 ns =      1073 µs =      1.07 ms =  0.0011 s

                []Struct Max     18764440 ns =     18764 µs =     18.76 ms =  0.0188 s
                []Struct Min       734005 ns =       734 µs =      0.73 ms =  0.0007 s
              []Struct Total  10900532318 ns =  10900532 µs =  10900.53 ms = 10.9005 s
                []Struct Avg      1090053 ns =      1090 µs =      1.09 ms =  0.0011 s

Done in  1m5.76463199s

```
