# TODO

## Google cloud connection issues
* frontend appengine app raises `panic: driver: bad connection` and won't retrieve database entry when deployed. Works when connecting to the same database from localhost.
* listener raises `Error 1045: Access denied for user`
    - Environment varibles do load
    - Able to ping and connect and ping the database from `pingtest` on the server.
    - Won't connect within the listener program
    - Connects when running locally
    - server running MySQL version 5.7; local 5.7.16
```
could not register possession:
[Error on RegisterPossession=
[Error on EndLastPossession#1=
 Error 1045: Access denied for user '<username>'@'<server_ip>' (using password: YES)]]"
```

## frontend interface
* markup leaderboard

## Database operations
* select top n possions of the month for leaderboard