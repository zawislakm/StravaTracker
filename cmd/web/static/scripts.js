// TODO after new table, sort it in order set by user
document.addEventListener("DOMContentLoaded", ()=>{
    const tableContainer = document.getElementById("table-container");
    const buttons = document.querySelectorAll(".btn-sort"); // Add dot for class selector
    console.log("STARTED", buttons.length);

    buttons.forEach((button) => {
        button.addEventListener("click", () => {
            const key = button.getAttribute("data-key"); // Match the attribute name in HTML
            const order = button.getAttribute("data-order");

            const headers = Array.from(document.querySelectorAll("th .header-content span"));
            const columnIndex = headers.findIndex(span => span.textContent === key) + 1; // +1 because of the Name column
            console.log(headers,columnIndex)

            const rows = Array.from(tableContainer.querySelectorAll("tbody tr"));

            const initialPositions = rows.map((row) => row.offsetTop);
            console.log(key, order, initialPositions);

            rows.sort((a,b) =>{
                const aValue = parseFloat(a.cells[columnIndex].textContent.trim().replace(/[^0-9]/g, ""));
                const bValue = parseFloat(b.cells[columnIndex].textContent.trim().replace(/[^0-9]/g, ""));
                console.log(aValue, bValue);
                if(order === "asc"){
                    return aValue > bValue ? -1 : 1;
                }else{
                    return aValue < bValue ? -1 : 1;
                }
            })

            rows.forEach((row,index) => {
                const newPosition = initialPositions[index]
                const currentPosition = row.offsetTop;

                setTimeout(() => {
                    row.style.transform = `translateY(${newPosition - currentPosition}px)`;
                    row.classList.add("moving")

                }, index * 100);
            })

            void tableContainer.offsetHeight;

        });
    });
});