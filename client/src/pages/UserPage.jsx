import {useRoute} from "wouter"

const UserPage = ()=> {
	const [match, params] = useRoute("/users/:username")
	// ensure that the user has the correct JWT
	return <h1>Hello {params.username} </h1>
}

export default UserPage;
