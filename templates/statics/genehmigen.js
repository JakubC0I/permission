function docReady(fn) {
    // see if DOM is already available
    if (document.readyState === "complete" || document.readyState === "interactive") {
        // call on next available tick
        setTimeout(fn, 1);
    } else {
        document.addEventListener("DOMContentLoaded", fn);
    }
}
async function put() {
    const res = await fetch(location.pathname, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        }
    })
    res.json().then(e => {
        console.log(e);
    })
}

docReady(put)