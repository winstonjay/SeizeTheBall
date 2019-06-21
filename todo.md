# TODO

## Google cloud connection issues
* frontend appengine app raises `panic: driver: bad connection` and won't retrieve database entry when deployed. Works when connecting to the same database from localhost.
* listener raises `mysql Error 1045: Access denied for user`
    - Are able to connect to the server from the computue engine but not within the listener program
    - Seems database environment varibles are not set (twitter ones are however)

## frontend interface
* markup leaderboard

## Database operations
* select top n possions of the month for leaderboard