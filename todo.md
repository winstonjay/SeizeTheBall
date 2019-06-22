# TODO

## Google cloud connection issues
* frontend appengine app raises `panic: driver: bad connection` and won't retrieve database entry when deployed. Works when connecting to the same database from localhost.
* listener raises `Error 1045: Access denied for user`
    - Environment varibles do load
    - Able to ping and connect and ping the database from `pingtest` on the server.
    - Won't connect within the listener program
    - Connects when running locally
    - server running MySQL version 5.7; local 5.7.16
    - Error is raised on `select` operation. It connects (or maybe lazy connects and dosen't fail until we try and access something)
    - `grant all privileges on *.* to '<username>'@'<ip>';` gives `ERROR 1045 (28000): Access denied for user`
    - Changed port, changed charset didnt make any difference.
    - switched to google mysql database
```
could not register possession:
[Error on RegisterPossession=
[Error on EndLastPossession#1=
 Error 1045: Access denied for user '<username>'@'<server_ip>' (using password: YES)]]"
```

tests now failing.
```
go test ./... -cover -v
?       _/Users/laurie/Documents/_code/go/src/github.com/winstonjay/seizeTheBall/frontend       [no test files]
?       _/Users/laurie/Documents/_code/go/src/github.com/winstonjay/seizeTheBall/listener       [no test files]
?       _/Users/laurie/Documents/_code/go/src/github.com/winstonjay/seizeTheBall/logger [no test files]
=== RUN   TestConnection
--- PASS: TestConnection (0.44s)
=== RUN   TestRegisterPossession
--- FAIL: TestRegisterPossession (7.88s)
        model_test.go:65: Possession.End == nil at test 0
        model_test.go:65: Possession.End == nil at test 1
        model_test.go:65: Possession.End == nil at test 2
        model_test.go:65: Possession.End == nil at test 3
        model_test.go:65: Possession.End == nil at test 4
=== RUN   TestCurrentPossession
--- PASS: TestCurrentPossession (8.11s)
=== RUN   TestCreatePossesssion
--- PASS: TestCreatePossesssion (4.70s)
=== RUN   TestGetAllPossesssions
--- PASS: TestGetAllPossesssions (4.92s)
=== RUN   TestEndLastPoessession
--- FAIL: TestEndLastPoessession (8.05s)
        model_test.go:155: Possession.End == nil at test 0
        model_test.go:155: Possession.End == nil at test 1
        model_test.go:155: Possession.End == nil at test 2
        model_test.go:155: Possession.End == nil at test 3
        model_test.go:155: Possession.End == nil at test 4
        model_test.go:155: Possession.End == nil at test 5
=== RUN   TestCreateUser
--- PASS: TestCreateUser (1.95s)
=== RUN   TestGetUserID
--- PASS: TestGetUserID (3.33s)
=== RUN   TestGetOrCreateUser
--- PASS: TestGetOrCreateUser (4.83s)
FAIL
coverage: 67.3% of statements
FAIL    _/Users/laurie/Documents/_code/go/src/github.com/winstonjay/seizeTheBall/model  44.215s
?       _/Users/laurie/Documents/_code/go/src/github.com/winstonjay/seizeTheBall/ping   [no test files]
make: *** [test] Error 1
```

## frontend interface
* markup leaderboard

## Database operations
* select top n possions of the month for leaderboard