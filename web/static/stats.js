'use strict';

let stats = {
    today: undefined,
    week: undefined,
    month: undefined
}

const page = {

    xhr: new XMLHttpRequest(),

    init() {
        window.addEventListener('load', (event) => this.renderPage(event))
    }, 

    renderPage(_event) {
        let path = location.pathname.split('/');
        let shrt = path[path.length-1];
        this.xhr.open("GET", `/redirects/${shrt}`, true);
        console.log(location.pathname.split("/")[-1]);
        this.xhr.onreadystatechange = function () { // Call a function when the state changes.
            if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
                console.log(page.xhr.response);
                stats = JSON.parse(page.xhr.response);
                document.getElementById('stats').style.display = '';
                document.getElementById('today').innerText = stats.today;
                document.getElementById('week').innerText = stats.week;
                document.getElementById('month').innerText = stats.month;
                let link = document.getElementById('shrtn-h');
                link.innerText = `/${shrt}`;
                link.setAttribute('href',`/${shrt}`);
            }
        }
        this.xhr.send(null);
    }

}

page.init()
