document.addEventListener('DOMContentLoaded', function () {
    const navLinks = document.querySelectorAll('nav a');

    navLinks.forEach(link => {
        link.addEventListener('click', function (event) {
            navLinks.forEach(link => link.classList.remove('active'));
            event.target.classList.add('active');
        });
    });
});

document.addEventListener('DOMContentLoaded', function () {
    const navLinks = document.querySelectorAll('nav a');

    navLinks.forEach(link => {
        link.addEventListener('click', function (event) {
            navLinks.forEach(link => link.classList.remove('active'));
            event.target.classList.add('active');
        });
    });

    const sidebarToggle = document.getElementById('sidebarToggle');
    const sidebar = document.getElementById('sidebar');

    sidebarToggle.addEventListener('click', function () {
        sidebar.classList.toggle('-translate-x-full');
    });
});

document.addEventListener('DOMContentLoaded', function () {
    const navLinks = document.querySelectorAll('nav a');

    navLinks.forEach(link => {
        link.addEventListener('click', function (event) {
            navLinks.forEach(link => link.classList.remove('active'));
            event.target.classList.add('active');
        });
    });

    const sidebarToggle = document.getElementById('sidebarToggle');
    const sidebar = document.getElementById('sidebar');

    sidebarToggle.addEventListener('click', function () {
        sidebar.classList.toggle('-translate-x-full');
    });

    const menuTable = document.getElementById("menuTable");
    const filterInput = document.getElementById("filterInput");
    const categorySelect = document.getElementById("categorySelect");
    const headers = menuTable.querySelectorAll("th");
    let ascending = true;

    // Filter function
    filterInput.addEventListener("input", function () {
        const filterValue = filterInput.value.toLowerCase();
        const rows = menuTable.querySelectorAll("tbody tr");
        rows.forEach(row => {
            const nameCell = row.querySelector("td:nth-child(2)");
            const nameText = nameCell.textContent.toLowerCase();
            row.style.display = nameText.includes(filterValue) ? "" : "none";
        });
    });

    // Category filter function
    categorySelect.addEventListener("change", function () {
        const selectedCategory = categorySelect.value;
        const rows = menuTable.querySelectorAll("tbody tr");
        rows.forEach(row => {
            const categoryCells = row.querySelectorAll("td:nth-child(4) span");
            const hasCategory = Array.from(categoryCells).some(span => span.textContent === selectedCategory);
            row.style.display = selectedCategory === "All" || hasCategory ? "" : "none";
        });
    });

    // Sorting function
    headers.forEach((header, index) => {
        header.addEventListener("click", function () {
            const rows = Array.from(menuTable.querySelectorAll("tbody tr"));
            const sortedRows = rows.sort((a, b) => {
                const aText = a.querySelector(`td:nth-child(${index + 1})`).textContent;
                const bText = b.querySelector(`td:nth-child(${index + 1})`).textContent;
                return ascending ? aText.localeCompare(bText) : bText.localeCompare(aText);
            });
            ascending = !ascending;
            sortedRows.forEach(row => menuTable.querySelector("tbody").appendChild(row));
        });
    });
});
document.addEventListener('DOMContentLoaded', function () {
    const navLinks = document.querySelectorAll('nav a');

    navLinks.forEach(link => {
        link.addEventListener('click', function (event) {
            navLinks.forEach(link => link.classList.remove('active'));
            event.target.classList.add('active');
        });
    });

    const sidebarToggle = document.getElementById('sidebarToggle');
    const sidebar = document.getElementById('sidebar');

    sidebarToggle.addEventListener('click', function () {
        sidebar.classList.toggle('-translate-x-full');
    });

    const menuTable = document.getElementById("menuTable");
    const filterInput = document.getElementById("filterInput");
    const categorySelect = document.getElementById("categorySelect");
    const headers = menuTable.querySelectorAll("th");
    let ascending = true;

    // Filter function
    filterInput.addEventListener("input", function () {
        const filterValue = filterInput.value.toLowerCase();
        const rows = menuTable.querySelectorAll("tbody tr");
        rows.forEach(row => {
            const nameCell = row.querySelector("td:nth-child(2)");
            const nameText = nameCell.textContent.toLowerCase();
            row.style.display = nameText.includes(filterValue) ? "" : "none";
        });
    });

    // Category filter function
    categorySelect.addEventListener("change", function () {
        const selectedCategory = categorySelect.value;
        const rows = menuTable.querySelectorAll("tbody tr");
        rows.forEach(row => {
            const categoryCells = row.querySelectorAll("td:nth-child(4) span");
            const hasCategory = Array.from(categoryCells).some(span => span.textContent === selectedCategory);
            row.style.display = selectedCategory === "All" || hasCategory ? "" : "none";
        });
    });

    // Sorting function
    headers.forEach((header, index) => {
        header.addEventListener("click", function () {
            const rows = Array.from(menuTable.querySelectorAll("tbody tr"));
            const sortedRows = rows.sort((a, b) => {
                const aText = a.querySelector(`td:nth-child(${index + 1})`).textContent;
                const bText = b.querySelector(`td:nth-child(${index + 1})`).textContent;
                return ascending ? aText.localeCompare(bText) : bText.localeCompare(aText);
            });
            ascending = !ascending;
            sortedRows.forEach(row => menuTable.querySelector("tbody").appendChild(row));
        });
    });
});