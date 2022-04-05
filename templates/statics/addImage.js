const form = document.getElementById("form");
const image = form.file
form.addEventListener("submit", async (e) => {
    e.preventDefault()
    const description = form.description.value
    const title = form.title.value
    const article = form.article.value
    const res = await fetch("/addImage", {
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ "image": files, description, title, article}),
        method: "POST",
    })
    files = []
    res.json().then((data) => {
        const {images} = data
        console.log(images);
        images.forEach(element => {
                document.getElementById("imagesDiv").innerHTML += "<img src="+element+"><br>"
        });
    })

})
let files = [];
image.addEventListener('change', async (e) => {
    document.getElementById("imagesDiv").innerHTML = ""
    const fileList = form.file.files
    const fileNames = [];
    for (let index = 0; index < fileList.length; index++) {
        const element = fileList[index];
        const reader = new FileReader();
        reader.onload = (async () => {
            files.push(reader.result)
        })
        reader.readAsDataURL(element)
        fileNames.push(element.name)
    };
}
)