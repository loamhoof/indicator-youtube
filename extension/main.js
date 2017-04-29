;(() => {
    const TARGET = 'http://127.0.0.1:12346'

    setInterval(() => {
        chrome.tabs.query({url: 'https://www.youtube.com/watch*'}, (tabs) => {
            let xhr = new XMLHttpRequest();

            for (let tab of tabs) {
                if (tab.audible) {
                    chrome.tabs.executeScript(tab.id, {code: '(' + (() => {
                        let video = document.getElementsByTagName('video');

                        if (!video.length) {
                            return [0, 0];
                        }
                        return [video[0].currentTime || 0, video[0].duration || 0];
                    }) + ')();'}, ([[current, duration]]) => {
                        xhr.open('GET', `${TARGET}/play/${encodeURIComponent(tab.title)}/${current}/${duration}`);
                        xhr.send();
                    });

                    return;
                }
            }

            let path = tabs.map((tab) => encodeURIComponent(tab.title)).join('/');

            xhr.open('GET', `${TARGET}/pause/${path}`);
            xhr.send();
        })
    }, 1000);
})();
