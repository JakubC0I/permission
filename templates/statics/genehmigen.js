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
    await fetch(location.pathname, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        }
    })
    console.log(res);
}

docReady(put)