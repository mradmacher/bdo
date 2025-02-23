# Development
1. Create `.env` file with the following content:
 ```
 export BDO_DB_URI="./db/bdo.sqlite"
 export GOOGLE_MAPS_API_KEY="<here-goes-the-key>"
```
2. Run `npm install`
3. Open the database `sqlite3 db/bdo.sqlite` and create the schema `.read db/sql/create.sql`
4. Seed the database with `go run cmd/dbseed/main.go`
5. Start the web server `bin/start`
