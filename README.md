<<<<<<< HEAD
#My first Api restful

components:

1. Go:
	-Chi
	-Cockroach - Posgresql
	-net/http
	-whois-go
	-other things.

2. HTML/Boostrap 4/Javascript
	-Vue js

Compiling:	

1. Creating one-clustring of one node:

	$ cockroach start-single-node --insecure --store=json-test --listen-addr=localhost:26257 --http-addr=localhost:8080 --background
	$ cockroach sql --insecure
	
	\>CREATE USER IF NOT EXISTS maxroach;
	\>CREATE DATABASE bank;
	\>GRANT ALL ON DATABASE bank TO maxroach;
	\> \q

For entering to database from terminal:

	$ cockroach sql --insecure -e --database=bank

2. compling golang from terminal and getting up the localhost server:
	
	$ go1.10.7 run api_rest.go

and all right!! watch at the browser, and put por instance:

	"http://localhost:8000/=example.com"

3. Finally in other port open "index.html", and experiment with my app.
=======
# App_api_rest
>>>>>>> 7ab50243970bad0744dcf16e3658a3a08e2fd502
