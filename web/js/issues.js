function test(url) {

    console.log("Testing URL: " + url);

    try {
        fetch(url)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok ' + response.statusText);
                }
                return response.json();
            })
            .then(data => {
                console.log(data);
                console.log(data[1]);
                console.log(data[1].title);
            })
            .catch(error => {
                console.error('There has been a problem with your fetch operation:', error);
            });
    } catch (e) {
        console.error(e);
    }
}
