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
		body.forEach(obj => {
			const newRow = table.insertRow(0)
			newRow.insertCell(0).innerHTML = obj.site
			newRow.insertCell(1).innerHTML = obj.site_user_name
			const password = newRow.insertCell(2)
			password.style.color = "transparent"
			password.innerHTML = obj.password
		})
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
