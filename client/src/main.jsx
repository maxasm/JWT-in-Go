import React from "react"
import ReactDOM from "react-dom/client"
import Login from "./pages/Login"

const root_div = document.getElementById("root")
document.title = "Hello"

ReactDOM.createRoot(root_div).render(
	<React.StrictMode>
		<Login/>	
	</React.StrictMode>
)
