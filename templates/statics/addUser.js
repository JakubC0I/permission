const form = document.getElementById("form")


form.addEventListener("submit", async e => {
    e.preventDefault()
    const email = form.email.value
    const password = form.password.value
    const genehmiger = form.genehmiger.value
    let res = await fetch("/addUser",{
        method: "POST",
        body: JSON.stringify({ email, password, genehmiger }),
        headers: {
            "Content-Type": "application/json"
        }
    })

    res.json().then( result => {
        console.log(result);
    })
})