# Seize the Ball!

**WORK IN PROGRESS**

A virtual ball game to play on Twitter. Seize the ball by tweeting "@seizeTheBall I have the ball" and see how long you manage keep hold of it.

## Environment Setup
You will need a MySQL server and a [Twitter developer](https://developer.twitter.com/content/developer-twitter/en.html) account and an app created to generate the credentials required for this application.

### To run locally

**Set enviroment varibles**; create a `.env` file and fill in the credientials as described in the `.env.example` file to allow access to your MySQL database and Twitter app. These need to be exported to be used in your shell session; you can export all of them in one line with:

```
$ for line in $(cat .env); do export $line; done
```

**Setup database**: Using MySQL 5.6 or later create the database schema.

```
$ make localdb
```

Also need to install the go MySQL driver:
```
$ go get "github.com/go-sql-driver/mysql"
```
