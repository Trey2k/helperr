/* eslint-disable no-undef */
let timeout;

const password = document.getElementById('password');
const label = document.getElementById('strengthLabel');
const badge = document.getElementById('strengthBadge');
const button = document.querySelector('input[type="submit"]');

const regexStrong = new RegExp('(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[^A-Za-z0-9])(?=.{8,})');
const regexMedium = new RegExp('((?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[^A-Za-z0-9])(?=.{6,}))|((?=.*[a-z])(?=.*[A-Z])(?=.*[^A-Za-z0-9])(?=.{8,}))');

function strengthCheck(passwordParam) {
	if(regexStrong.test(passwordParam)) {
		button.removeAttribute('disabled');
		badge.style.color = 'green';
		badge.textContent = 'Strong';
	}
	else if (regexMedium.test(passwordParam)) {
		button.setAttribute('disabled', 'disabled');
		badge.style.color = 'orange';
		badge.textContent = 'Medium';
	}
	else {
		button.setAttribute('disabled', 'disabled');
		badge.style.color = 'red';
		badge.textContent = 'Weak';
	}
}

password.addEventListener('input', () => {
	label.style.display = 'block';
	badge.style.display = 'inline-block';
	clearTimeout(timeout);

	timeout = setTimeout(() => strengthCheck(password.value), 500);

	if(password.value.length !== 0) {
		label.style.display = 'block';
		badge.style.display != 'inline-block';
	}
	else {
		button.setAttribute('disabled', 'disabled');
		label.style.display = 'none';
		badge.style.display = 'none';
	}
});
