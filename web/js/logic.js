(function() {

    oldTags = [];

    taggedKvTemplate = document.getElementById("tagged-key-value-template");
    taggedKvContainer = document.getElementById("tagged-key-values-container");
    search = document.getElementById("search-input");
    search.addEventListener("keydown", (event) => {
        if ((event.key === "Enter") || (event.key === " ")) {
            tags = search.value === "*" ? [] : search.value.split(" ");

            if (tags.toString() === oldTags.toString()) {
                return;
            }

            console.log("searching for: ", tags);

            // clear existing entries
            while (taggedKvContainer.firstChild) {
                taggedKvContainer.removeChild(taggedKvContainer.firstChild);
            }

            // search key
            fetch(`api/keys/${search.value}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            })
            .then(response => response.json())
            .then(data => {
                console.log("received data: ", data);

                data.forEach((item) => {
                    let clone = taggedKvTemplate.content.cloneNode(true);
                    clone.querySelector(".tkv-title").innerText = item.key;
                    clone.querySelector(".tkv-value").innerText = item.value;

                    flatTags = (item.tags || []).join(" ");
                    clone.querySelector(".tkv-tags").innerText = "Tags: " + flatTags;

                    taggedKvContainer.appendChild(clone);
                });

            })
            .catch((error) => {
                console.error('Error:', error);
            });

            // search tags
            fetch(`api/keys?tags=${tags.toString()}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            })
            .then(response => response.json())
            .then(data => {
                console.log("received data: ", data);

                data.forEach((item) => {
                    let clone = taggedKvTemplate.content.cloneNode(true);
                    clone.querySelector(".tkv-title").innerText = item.key;
                    clone.querySelector(".tkv-value").innerText = item.value;

                    flatTags = (item.tags || []).join(" ");
                    clone.querySelector(".tkv-tags").innerText = "Tags: " + flatTags;

                    taggedKvContainer.appendChild(clone);
                });

            })
            .catch((error) => {
                console.error('Error:', error);
            });

            oldTags = tags;
        }
    });
})()
