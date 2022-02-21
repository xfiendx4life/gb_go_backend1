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

let stats = {
    day: undefined,
    week: undefined,
    month: undefined
}

const page = {

    u: user,
    url: url,
    stats: stats,
    xhr: new XMLHttpRequest(),
    init() {
        window.addEventListener('load', (event) => this.renderPage(event))
        window.addEventListener('click', (event) => this.ClickEvent(event));
    },
    jwt: undefined,

    getCookie(key) {
        const c = document.cookie
        let res = 0
        if (c != '') {
            res = c.split('; ')
                .find(row => row.startsWith(`${key}=`))
                .split('=')[1];
        }

        return c
    },

    renderPage(_event) {
        if (document.cookie !== "") {
            const jwt = this.getCookie('jwt');
            const id = +this.getCookie('id');
            if (jwt !== '' && id != '' && id != 0) {
                this.jwt = jwt;
                this.u.id = id;
                this.getToPhase2(false);
            } else {

                document.getElementById('reg-form').style.display = '';
            }
        } else {
            document.getElementById('reg-form').style.display = '';
        }

    },


    login() {
        let target = `/user/login?name=${this.u.name}&password=${this.u.password}`
        console.log(target);
        this.xhr.open('GET', target, true);
        this.xhr.onload = function () {
            if (page.jwt === undefined || page.jwt === '"message":"Unauthorized"') {
                page.jwt = page.xhr.response.slice(1, -2);
                if (page.jwt !== '"message":"Unauthorized"') {
                    let f = document.getElementById('reg-form');
                    if (f) {
                        f.remove();
                    };
                    document.getElementById('shrtn-form').style.display = '';
                } else {
                    page.sendMessage('Login attempt failed, try again or register');
                }
                document.cookie = `jwt=${page.jwt}`;
            }

            if (page.id === undefined && page.u.id !== undefined) {
                page.id = page.u.id;
                document.cookie = `id=${page.u.id}`
            }
            if (page.u.id === undefined) {
                page.getUserId();
            }

            console.log(page.jwt);
        }
        this.xhr.send(null);
    },

    getUserId() {
        let target = `/user?name=${this.u.name}`
        console.log(target);
        this.xhr.open('GET', target, true);
        this.xhr.setRequestHeader("Authorization", "Bearer " + this.jwt);
        this.xhr.onload = function () {
            if (page.u.id === undefined) {
                page.id = page.xhr.response.slice(0, -1);
                page.u.id = page.id;
                document.cookie = `id=${page.id}`;
            }
        }
        this.xhr.send(null);
    },

    getToPhase2(shallLogin) {
        let f = document.getElementById('reg-form');
        if (shallLogin) {
            this.login();
        }
        if (+this.getCookie('id') !== 0 && this.getCookie('jwt') !== '"message":"Unauthorized"') {
            f.remove();
            document.getElementById('shrtn-form').style.display = '';
        }

    },

    sendMessage(text) {
        const alert = document.getElementById('alert');
        document.getElementById('alert-text').innerHTML = text;
        alert.style.display = '';
    },

    ClickEvent(event) {
        const target = event.target.className;
        if (event.target.tagName === 'BUTTON') {
            if (target === 'btn btn-secondary register-btn') {
                this.getDataFromForm();
                this.xhr.open("POST", `/user/create`, true);
                this.xhr.setRequestHeader("Content-Type", "application/json");
                this.xhr.onreadystatechange = function () { // Call a function when the state changes.
                    if (this.readyState === XMLHttpRequest.DONE && this.status === 201) {
                        console.log(page.xhr.response)
                        page.u.id = JSON.parse(page.xhr.response).id
                        page.getToPhase2(true);
                        document.getElementById('shrtn-form').style.display = '';

                    }
                }
                this.xhr.send(JSON.stringify(this.u));

            } else if (target === 'btn btn-primary shorten') {
                this.url.rawurl = this.getUrlData();
                this.url.userid = +this.u.id;
                this.xhr.open("POST", `/url`, true);
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

                        let stats = document.getElementById('stats-link');
                        stats.setAttribute('href', `${window.location.host}/web/redirects/${page.url.shortened}`);
                        stats.innerText = `/web/redirects/${page.url.shortened}`;

                    }
                }
                this.xhr.send(JSON.stringify(this.url));
            } else if (target === 'btn btn-primary login-btn') {
                this.getDataFromForm();
                this.getToPhase2(true);
                // if (document.cookie != '' && +this.getCookie('id') !== 0 && this.getCookie('jwt') !== '"message":"Unauthorized"') {
                //     document.getElementById('shrtn-form').style.display = '';
                // }
            } else if (target === 'btn btn-link') {
                document.getElementById('stats').style.display = '';
                const url = document.getElementById('stats-link').innerText;
                this.xhr.open("GET", `/redirects${url}`, true);
                this.xhr.setRequestHeader("Content-Type", "application/json");
                this.xhr.setRequestHeader("Authorization", "Bearer " + this.jwt);
                this.xhr.onreadystatechange = function () { // Call a function when the state changes.
                    if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
                        console.log(page.xhr.response);
                        stats = JSON.parse(page.xhr.response);
                        document.getElementById('today').innerText = stats.day;
                        document.getElementById('week').innerText = stats.week;
                        document.getElementById('month').innerText = stats.month;
                    }
                }
                this.xhr.send(null);
            }

        }
    },

    getDataFromForm() {
        let login = document.getElementById('login').value;
        let password = document.getElementById('password').value;
        let email = document.getElementById('email').value;
        this.u.init(login, password, email);
        this.u.id = undefined;
    },

    getUrlData() {
        return document.getElementById('raw').value;

    }
}
page.init()