document.addEventListener('DOMContentLoaded', function () {
    const navLinks = document.querySelectorAll('nav a');

    navLinks.forEach(link => {
        link.addEventListener('click', function (event) {
            navLinks.forEach(link => link.classList.remove('active'));
            event.target.classList.add('active');
        });
    });
});