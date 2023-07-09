import Stack from "@mui/material/Stack"
import TextField from "@mui/material/TextField"
import Button from "@mui/material/Button"
import InputAdornment from "@mui/material/InputAdornment"
import Visibility from "@mui/icons-material/Visibility"
import VisibilityOff from "@mui/icons-material/VisibilityOff"
import IconButton from "@mui/material/IconButton"
import OutlinedInput from "@mui/material/OutlinedInput"
import FormControl from "@mui/material/FormControl"
import InputLabel from "@mui/material/InputLabel"
import Snackbar from "@mui/material/Snackbar"

// routing
import {useLocation} from "wouter"
// colors
import blue from "@mui/material/colors/blue"

import {useState} from "react"


const Login = ()=> {
// hide/show password
const [hide,updateHide] = useState(true)

// username and password
const [username,updateUsername] = useState("")
const [password,updatePassword] = useState("")

function handleUpdateUsername(e) {
	updateUsername(e.target.value)
}

function handleUpdatePassword(e) {
	updatePassword(e.target.value)
}

// routing
const [location,updateLocation] = useLocation()

// Login using username and password
async function handleLogin() {
	const response = await fetch('/api_login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
	});

	let msg = await response.text()
	msg && alert(msg)
	
	if (response.status == 200) {
		// navigate to the user page
		updateLocation(`/users/${username}`)
	}
}


// Signup using username and password
async function handleSignup() {
	const response = await fetch('/api_signup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
	});
	
	let msg = await response.text()
	msg && alert(msg)
}


return (
<>
	<Stack spacing={4} sx={{width: 400, margin:"100px auto", padding: "20px", border: `2px solid ${blue[400]}`,borderRadius:"4px"}}>
		<TextField
			onChange={handleUpdateUsername}
			value={username}
			label="Username"
			variant="outlined"/>
		<FormControl>
			<InputLabel>Password</InputLabel>
			<OutlinedInput
				onChange={handleUpdatePassword}
				value={password}
				variant="outlined"
				label="Password"
				type= { hide ? "password" : "text" }
				endAdornment={
					<InputAdornment position="end">
						<IconButton onClick={()=> {updateHide(!hide)}}>
							{hide ? <VisibilityOff /> : <Visibility />}
						</IconButton>
					</InputAdornment>}/>
		</FormControl>
			<Stack spacing={2} direction="row">
				<Button onClick={handleLogin} variant="outlined">Log In</Button>
				<Button onClick={handleSignup} variant="outlined">Sign Up</Button>
			</Stack>
	</Stack>
</>)}

export default Login;
