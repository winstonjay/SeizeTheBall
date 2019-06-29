# TODO

## Google cloud to MySQL connection issues
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
    - DONE

## frontend interface
* markup leaderboard

## Database operations
* select top n possions of the month for leaderboard