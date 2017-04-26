;(() => {
    const TARGET = 'http://127.0.0.1:12346'

    setInterval(() => {
        chrome.tabs.query({url: 'https://www.youtube.com/*'}, (tabs) => {
            let xhr = new XMLHttpRequest();

            for (let tab of tabs) {
                if (tab.audible) {
                    xhr.open('GET', `${TARGET}/${encodeURIComponent(tab.title)}`);
                    xhr.send();

                    return;
                }
            }

            xhr.open('GET', `${TARGET}/`);
            xhr.send();
        })
    }, 1000);
})();
