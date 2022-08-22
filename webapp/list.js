if (localStorage.getItem("password_manager_jwt_token") == null) {
	window.location.replace(window.location.href.replace("list.html", "index.html"));
}

const remoteUrl = "https://password-manager-b7jpyqffha-uc.a.run.app"
const localUrl = "http://localhost:8080"
const url = window.location.href.startsWith("https") ? remoteUrl : localUrl;
async function getPasswords() {
	const req = {}
	const resp = await fetch(url + "/secure/listPasswords", {
		method: "POST",
		headers: {
			'Content-Type': 'application/json',
			'Authorization': 'Bearer ' + localStorage.getItem("password_manager_jwt_token"),
		},
		body: JSON.stringify(req),
	})
	if (resp.status == 400 || resp.status == 401) {
		localStorage.removeItem("password_manager_jwt_token")
		window.location.replace(window.location.href.replace("list.html", "index.html"));
	} else if (resp.status == 200) {
		const table = document.getElementById("passwords_body");
		const body = await resp.json()
		body.forEach(obj => insertRow(obj, table))
	}
}

function insertRow(obj, table) {
	const newRow = table.insertRow(0)
	newRow.id = obj.site + obj.site_user_name
	newRow.insertCell(0).innerHTML = obj.site
	newRow.insertCell(1).innerHTML = obj.site_user_name
	const password = newRow.insertCell(2)
	password.style.color = "transparent"
	password.innerHTML = obj.password
	newRow.insertCell(3).innerHTML = `
				<button type="button" class="btn btn-outline-danger btn-sm" onclick="deletePassword('${newRow.id}')">Delete Password</button>
	`
}
function deleteRow(rowId) {
	const row = document.getElementById(rowId)
	row.parentNode.removeChild(row)
}

async function deletePassword(rowId) {
	const columns = document.getElementById(rowId).children
	const req = {
		site: columns[0].innerText,
		site_user_name: columns[1].innerText,
	}
	const reallyDelete = confirm(`Are you sure you want to delete password for site ${req.site} and username ${req.site_user_name}`)
	if (!reallyDelete) {
		return
	}
	const resp = await fetch(url + "/secure/deletePassword", {
		method: "POST",
		headers: {
			'Content-Type': 'application/json',
			'Authorization': 'Bearer ' + localStorage.getItem("password_manager_jwt_token"),
		},
		body: JSON.stringify(req),
	})
	if (resp.status == 200) {
		deleteRow(rowId)
	} else if (resp.status == 401) {
		localStorage.removeItem("password_manager_jwt_token")
		window.location.replace(window.location.href.replace("list.html", "index.html"));
	} else {
		const body = await resp.json()
		console.log(body)
	}
}

async function submit() {
	const site = document.getElementById('site').value
	const username = document.getElementById('site_username').value
	const password = document.getElementById('password').value
	const req = {
		site: site,
		site_user_name: username,
		password: password,
	}
	const resp = await fetch(url + "/secure/upsertPassword", {
		method: "POST",
		headers: {
			'Content-Type': 'application/json',
			'Authorization': 'Bearer ' + localStorage.getItem("password_manager_jwt_token"),
		},
		body: JSON.stringify(req),
	})
	if (resp.status == 200) {
		document.getElementById("potential_success").hidden = false
		document.getElementById("potential_error").hidden = true
		document.getElementById("potential_success").innerHTML = "The password for site " + site + " and username " + username + " updated"
		insertRow(req, document.getElementById("passwords_body"))
	} else if (resp.status == 400) {
		document.getElementById("potential_success").hidden = true
		document.getElementById("potential_error").hidden = false
		document.getElementById("potential_error").innerHTML = "Please fill out all of the fields before submitting."
	} else if (resp.status == 401) {
		localStorage.removeItem("password_manager_jwt_token")
		window.location.replace(window.location.href.replace("list.html", "index.html"));
	} else {
		const body = await resp.json()
		document.getElementById("potential_error").innerHTML = body.message
	}
}


$(document).ready(async function() {
	await getPasswords()
	$("#password_site_search").on("keyup", function() {
		var value = $(this).val().toLowerCase();
		$("#passwords_body tr").filter(function() {
			const columns = $(this).children()
			const search = columns[0].innerText + columns[1].innerText
			$(this).toggle(search.toLowerCase().indexOf(value) > -1)
		});
	});
});
