;(() => {
    const TARGET = 'http://127.0.0.1:12346'

    setInterval(() => {
        chrome.tabs.query({url: 'https://www.youtube.com/watch*'}, (tabs) => {
            // Only consider one tab, don't want the same mess as with YT
            if (!tabs.length) {
                return;
            }

            let tab = tabs[0];

            let paused = !tab.audible;

            chrome.tabs.executeScript(tab.id, {code: '(' + (() => {
                let video = document.getElementsByTagName('video');

                if (!video.length) {
                    return [0, 0];
                }
                return [video[0].currentTime || 0, video[0].duration || 0];
            }) + ')();'}, ([[current, duration]]) => {
                let xhr = new XMLHttpRequest();
                xhr.open('GET', `${TARGET}/play/${encodeURIComponent(tab.title)}/${current}/${duration}/${paused}`);
                xhr.send();
            });
        })
    }, 1000);
})();
