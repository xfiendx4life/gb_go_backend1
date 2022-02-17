'use strict';

const user = {
    id: undefined,
    name: null,
    password: null,
    email: null,

    init(name, password, email) {
        this.name = name;
        this.password = password;
        this.email = email;
    },

}

const page = {

    u: user,
    xhr: new XMLHttpRequest(),
    init() {
        window.addEventListener('click', (event) => this.ClickEvent(event));
    },
    jwt: undefined,

    login() {
        const target = new URL('http://localhost:8080/user/login');
        const params = new URLSearchParams();
        params.set('name', this.u.name);
        params.set('password', this.u.password);
        target.search = params.toString();
        console.log(target);
        this.xhr.open('GET', target, true);
        this.xhr.onload = function () {
            this.jwt = page.xhr.response;
            console.log(this.jwt);
        }
        this.xhr.send(null);
    },

    ClickEvent(event) {
        const target = event.target.className;
        if (event.target.tagName === 'BUTTON') {
            if (target === 'btn btn-primary register-btn') {
                this.getDataFromForm();
                this.xhr.open("POST", "http://localhost:8080/user/create", true);
                this.xhr.setRequestHeader("Content-Type", "application/json");
                this.xhr.onreadystatechange = function () { // Call a function when the state changes.
                    if (this.readyState === XMLHttpRequest.DONE && this.status === 201) {
                        console.log(page.xhr.response)
                        page.u.id = JSON.parse(page.xhr.response).id
                        let f = document.getElementById('reg-form')
                        f.remove()
                        page.login()
                        document.getElementById('reg-form').style.display = '';
                        
                    }
                }
                this.xhr.send(JSON.stringify(this.u));

            }
        }
    },

    getDataFromForm() {
        let login = document.getElementById('login').value;
        let password = document.getElementById('password').value;
        let email = document.getElementById('email').value;
        this.u.init(login, password, email);
    }
}
page.init()