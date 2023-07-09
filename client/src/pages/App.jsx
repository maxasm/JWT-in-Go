import {useEffect, useState} from "react"


const App = ()=> {
	const [user,updateUser] = useState()
	
	useEffect(() => {
		async function fetchUser() {
			let res;
			try {
				let res = await fetch("/api_app")
				let data = await res.text()
				console.log(data)
				updateUser(data)
				// check for 401
				if (res.status === 401) {
					updateUser(null)
				}
			} catch (e) {
				console.log(e)
			}
			
		}

		fetchUser()
	},[])
	
	return <h1> {user == null ? `Unauthorized` : `Welcome ${user}`} </h1>	
}

export default App;
