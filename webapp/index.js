if (localStorage.getItem("password_manager_jwt_token") != null) {
	window.location.replace(window.location.href.replace("index.html", "list.html"));
	const req = {}
	const resp = await fetch(url + "/secure", {
		method: "POST",
		headers: {
			'Content-Type': 'application/json',
			'Authorization': 'Bearer ' + localStorage.getItem("password_manager_jwt_token"),
		},
		body: JSON.stringify(req),
	})
	if (resp.status == 400 || resp.status == 401) {
		localStorage.removeItem("password_manager_jwt_token")
	} else if (resp.status == 200) {
		window.location.replace(window.location.href.replace("index.html", "list.html"));
	}
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
	const resp = await fetch(url + "/authenticate", {
		method: "POST",
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(req),
	})
	if (resp.status == 400 || resp.status == 401) {
		document.getElementById("potential_error").innerHTML = "Either the email or the password is incorrect."
		document.getElementById("potential_error").hidden = false
	} else if (resp.status == 200) {
		document.getElementById("potential_error").hidden = true
		const body = await resp.json()
		localStorage.setItem("password_manager_jwt_token", body.token)
		window.location.replace(window.location.href.replace("index.html", "list.html"));
	}
}

