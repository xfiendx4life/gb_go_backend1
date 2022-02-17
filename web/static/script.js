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

const url = {
    id: undefined,
    rawurl: null,
    shortened: undefined,
    userid: null
}

const page = {

    u: user,
    url: url,
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
            if (page.jwt === undefined) {
                page.jwt = page.xhr.response.slice(1, -2);
                console.log(this.jwt);
            }

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
                        document.getElementById('shrtn-form').style.display = '';

                    }
                }
                this.xhr.send(JSON.stringify(this.u));

            } else if (target === 'btn btn-primary shorten') {
                this.url.rawurl = this.getUrlData();
                this.url.userid = this.u.id;
                this.xhr.open("POST", "http://localhost:8080/url", true);
                this.xhr.setRequestHeader("Content-Type", "application/json");
                this.xhr.setRequestHeader("Authorization", "Bearer " + this.jwt);
                this.xhr.onreadystatechange = function () { // Call a function when the state changes.
                    if (this.readyState === XMLHttpRequest.DONE && this.status === 201) {
                        console.log(page.xhr.response);
                        page.url.shortened = JSON.parse(page.xhr.response).shortened;
                        document.getElementById('links-div').style.display = '';
                        let shrtLink = document.getElementById('shrtn-link');
                        shrtLink.setAttribute('href', '/' + page.url.shortened);
                        shrtLink.innerText = '/' + page.url.shortened;

                        let stats = document.getElementById('stats');
                        stats.setAttribute('href', '/redirects/' + page.url.shortened);
                        stats.innerText = '/redirects/' + page.url.shortened;

                    }
                }
                this.xhr.send(JSON.stringify(this.url));
            }
        }
    },

    getDataFromForm() {
        let login = document.getElementById('login').value;
        let password = document.getElementById('password').value;
        let email = document.getElementById('email').value;
        this.u.init(login, password, email);
    },

    getUrlData() {
        return document.getElementById('raw').value;

    }
}
page.init()
