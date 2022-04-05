const form1 = document.querySelectorAll(".form1")
const form2 = document.querySelectorAll(".form2")

form1.forEach(ele => {
    ele.addEventListener("submit", async e => {
        e.preventDefault()
        await fetch(ele.action, {
            method: "PUT"
        })
    })
})

form2.forEach( ele => {
    ele.addEventListener("submit", async e => {
        e.preventDefault()
        await fetch(ele.action, {
            method: "DELETE"
        })
    })
    
})