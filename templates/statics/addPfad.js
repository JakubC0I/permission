const form = document.getElementById("form")

form.addEventListener("submit", async(e) => {
    e.preventDefault()
    let ids = form.ids.value.split(" ")
    let pfads = form.pfads.value.split(" ")
    
    const res = await fetch("/addPfad", {
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ "_bids": ids, "data": pfads, 
        // "besteller": form.besteller.value
    }),
        method: "PUT",
    })
    res.json().then((data)=> {
        console.log(data)
    })
})