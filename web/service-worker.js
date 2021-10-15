const CACHE_NAME = "craterun-v1.1.1";
const CACHED_FILES = [
    "/crate.audio.worklet.js",
    "/crate.js",
    "/crate.pck",
    "/crate.png",
    "/crate.wasm",
    "/index.html",
    "/",
    "/favicon.png",
    "/pwa-icon.png",
    "/service-worker.js",
];

self.addEventListener("install", (event) => {
    event.waitUntil(
        (async () => {
            const cache = await caches.open(CACHE_NAME);
            await cache.addAll(CACHED_FILES);
        })()
    );
});

self.addEventListener("activate", (event) => {
    event.waitUntil(
        caches.keys().then((keyList) => {
            return Promise.all(
                keyList.map((key) => {
                    if (key === CACHE_NAME) {
                        return;
                    }

                    return caches.delete(key);
                })
            );
        })
    );
});

self.addEventListener("fetch", (event) => {
    event.respondWith(
        (async () => {
            const cached = await caches.match(event.request);

            if (cached) {
                return cached;
            }

            const response = await fetch(event.request);

            return response;
        })()
    );
});
