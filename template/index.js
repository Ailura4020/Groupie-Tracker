function handlerSubmit(event) {
    event.preventDefault();
    const searchQuery = document.querySelector(".js-search-input").value.trim();
    const searchResults = document.querySelector(".js-search-result");
    searchResults.innerHTML = "";

    const endpoint = `/search?q=${encodeURIComponent(searchQuery)}`;

    fetch(endpoint)
        .then((response) => {
            if (!response.ok) throw Error(response.statusText);
            return response.json();
        })
        .then((results) => {
            if (results.length === 0) {
                alert("No result found");
                return;
            }
            results.forEach((result) => {
                const url = `https://groupietrackers.herokuapp.com/api/id${result.id}`;
                searchResults.insertAdjacentHTML(
                    "beforeend",
                    `
                        <div class="MainBlock">
                            <div class="GroupInfo">
                                <p class="GroupName">${result.name}</p>
                                <p class="GroupStyle">Style: ${result.style}</p>
                                <a href="${url}/artists?id=${result.id}">
                                    <img class="GroupImg" src="${result.image}">
                                    ${result.name}
                                </a>
                            </div>
                        </div>
                    `
                );
            });
        })
        .catch((err) => {
            console.log(err);
            alert("Echec de la recherche");
        });
}

document.querySelector(".js-search-form").addEventListener("submit", handlerSubmit);