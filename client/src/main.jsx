import React from "react"
import ReactDOM from "react-dom/client"
import Login from "./pages/Login"
import UserPage from "./pages/UserPage"
import App from "./pages/App"

// routing
import {Route, Switch} from "wouter"

const root_div = document.getElementById("root")
document.title = "Hello"

ReactDOM.createRoot(root_div).render(
	<React.StrictMode>
		<Switch>
			<Route path="/userlogin">
				<Login/>
			</Route>
			<Route path="/users/:name">
				<UserPage/>
			</Route>
			<Route path="/app">
				<App/>
			</Route>
			<Route>
				<h1> 404 Page Not Found </h1>
			</Route>
		</Switch>
	</React.StrictMode>
)
