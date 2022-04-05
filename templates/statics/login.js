const form = document.getElementById("form")


form.addEventListener("submit", async e => {
    e.preventDefault()
    const email = form.email.value
    const password = form.password.value
    let res = await fetch("/login",{
        method: "POST",
        body: JSON.stringify({ email, password }),
        headers: {
            "Content-Type": "application/json"
        }
    })

    res.json().then( result => {
        console.log(result);
    })
})