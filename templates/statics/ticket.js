const form = document.getElementById("form")
const sForm = document.getElementById("sForm")

if (form) {
    form.addEventListener("submit", async(e) => {
        e.preventDefault()
        const comment = form.commentField.value
        let res = await fetch("/addComment", {
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({comment}),
            method: "POST"
        })
        res.json().then((result) => {
            console.log(result);
        })
    })
} else {
    sForm.addEventListener("submit", async(e) => {
        e.preventDefault()
        const searchbar = sForm.searchbar.value
        let res = await fetch("/search", {
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({searchbar}),
            method: "POST"
        })
        res.json().then((result) => {
            document.getElementById("incidents").innerHTML = ``
            result.data.forEach(ele => {
                console.log(ele);
                document.getElementById("incidents").innerHTML += `<a href="http://localhost:4000/ticket/${ele[2]}">${ele[1]}</a><br><p>${ele[0]}</p>`
            })
        })
    })
}