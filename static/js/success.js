// display the name of the student that signed in
// and the current time

const urlParams = new URLSearchParams(window.location.search);
const firstName = urlParams.get('firstName');

document.querySelector('.welcome-title').innerHTML = `Welcome ${firstName}!`;
document.querySelector('.time').innerHTML = new Date().toLocaleString();

// after successful signin, the success page is displayed
// after 3 seconds, the page is redirected to the home page

setTimeout(function(){
    window.location.href = "/";
}, 3000);
