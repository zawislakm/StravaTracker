// TODO make it responsive
document.addEventListener("DOMContentLoaded", () => {
    const tableContainer = document.getElementById("table-container");
    const buttons = document.querySelectorAll(".btn-sort");
    const loadingOverlay = document.getElementById("loading-overlay");
    const MIN_LOADING_TIME = 400;

    let loadingStartTime = 0;

    document.body.addEventListener('htmx:beforeSwap', function (event) {
        // Check if the swap target is our table body
        if (event.detail.target.id === 'data-table-body') {
            // Record the start time, show loading spinner
            loadingStartTime = Date.now();
            loadingOverlay.classList.remove('hidden');

        }
    });

    document.body.addEventListener('htmx:afterSwap', function (event) {
        // Check if the swap target was our table body
        if (event.detail.target.id === 'data-table-body') {
            sortTable();

            // Calculate how long the loading has been visible, hide loading spinner after ensuring minimum display time
            const elapsedTime = Date.now() - loadingStartTime;
            const remainingTime = Math.max(0, MIN_LOADING_TIME - elapsedTime);
            setTimeout(() => {
                loadingOverlay.classList.add('hidden');
            }, remainingTime);
        }
    });

    buttons.forEach((button) => {
        button.addEventListener("click", () => {
            const key = button.getAttribute("data-key");
            const order = button.getAttribute("data-order");

            tableContainer.setAttribute("data-key", key);
            tableContainer.setAttribute("data-order", order);

            sortTable();
        });
    });
});


function sortTable(storedPositions = null) {
    const tableContainer = document.getElementById("table-container");
    const key = tableContainer.getAttribute("data-key");
    const order = tableContainer.getAttribute("data-order");

    const headers = Array.from(document.querySelectorAll("th .header-content span"));
    const columnIndex = headers.findIndex(span => span.textContent === key) + 1; // +1 because of the Name column

    const rows = Array.from(tableContainer.querySelectorAll("tbody tr"));

    const initialPositions = storedPositions || rows.map((row) => row.offsetTop);

    rows.sort((a, b) => {
        const aValue = parseFloat(a.cells[columnIndex].textContent.trim().replace(/[^0-9.]/g, ""));
        const bValue = parseFloat(b.cells[columnIndex].textContent.trim().replace(/[^0-9.]/g, ""));

        if (order === "asc") {
            return aValue > bValue ? -1 : 1;
        } else {
            return aValue < bValue ? -1 : 1;
        }
    });


    rows.forEach((row, index) => {
        const newPosition = initialPositions[index];
        const currentPosition = row.offsetTop;

        setTimeout(() => {
            row.style.transform = `translateY(${newPosition - currentPosition}px)`;
            row.classList.add("moving");
        }, index * 100);
    });

    void tableContainer.offsetHeight;
}

