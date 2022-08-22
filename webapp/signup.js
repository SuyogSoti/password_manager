if (localStorage.getItem("password_manager_jwt_token") != null) {
	window.location.replace(window.location.href.replace("signup.html", "list.html"));
}

const remoteUrl = "https://password-manager-b7jpyqffha-uc.a.run.app"
const localUrl = "http://localhost:8080"
const url = window.location.href.startsWith("https") ? remoteUrl : localUrl;
async function submit() {
	const email = document.getElementById('email').value
	const password = document.getElementById('password').value
	const req = {
		email: email,
		password: password,
	}
	const resp = await fetch(url + "/createUser", {
		method: "POST",
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(req),
	})
	if (resp.status == 400 || resp.status == 401) {
		const body = await resp.json()
		document.getElementById("potential_error").hidden = false
		document.getElementById("potential_error").innerHTML = body.message
	} else if (resp.status == 200) {
		document.getElementById("potential_error").hidden = true
		const body = await resp.json()
		localStorage.setItem("password_manager_jwt_token", body.token)
		window.location.replace(window.location.href.replace("signup.html", "list.html"));
	}
}


