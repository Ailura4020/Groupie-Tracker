let sorteNames = names.sort();
console.log(sorteNames);


//reference
let input = document.getElementById("input");

// execute function
input.addEventListener("keyup", (e)=>{
    removeElement();

for(let i of sorteNames){

    if(
        i.toLowerCase().startsWith(input.value.toLowerCase()) && 
        input.value !=""
    ) {
        let listItem = document.createElement("li")
        listItem.classList.add("list-items");
        listItem.style.cursor="pointer";
        listItem.setAttribute("onclick","displayNames('" + i + "')");
        let word = "<b>" + i.substring(0,input.value.length) + "</b>";
        word+= i.substring(input.value.length);
    
        listItem.innerHTML = word;
        document.querySelector(".list").appendChild(listItem);
    }
    }
});
function displayNames(value){
    input.value = value;
    removeElement();
}
function removeElement(){
    //clear all the item
    let items = document.querySelectorAll(".list-items");
    items.forEach((item) =>{
        item.remove()
        ;})
}