# go-websocket-chat
# Usage
- Run the server(websocket gateway) <br />
       `go run server/main.go` <br />
- Run authentication service (required for the app to run) <br />
       `go run services/auth/main.go` <br />
- Run sync service (required for the app to run) <br />
       `go run services/sync/main.go` <br />
- Run the React client: <br />
       `cd client/` <br />
       `npm install` <br />
       `npm start`
- Open [http://localhost:3000](http://localhost:3000) to use the app. The auth service will run on localhost:9090, and sync will run on localhost:9091