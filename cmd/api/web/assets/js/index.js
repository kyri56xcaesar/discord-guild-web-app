const sideBarAffectedElements = document.querySelectorAll('.sidebar-affected');
// console.log(sideBarAffectedElements)

const mainSideBar = document.getElementById("main-siderbar");
const navbtn = document.getElementById("navbtn");
const togglebtn = document.getElementById("togglebtn")
const guild_name = '$DADS';
const default_icon = "assets/icons/mozart_draw.jpg";


function openNav() {
    mainSideBar.style.width = "250px";
    navbtn.style.marginLeft = "20px";
    togglebtn.innerText = "→";

    for (let element of sideBarAffectedElements) {
        element.style.marginLeft = '250px'; 
    }
}

function closeNav() {
    mainSideBar.style.width = "0";
    navbtn.style.marginLeft = "0";
    togglebtn.innerText = "←";

    for (let element of sideBarAffectedElements) {
        element.style.marginLeft = ''; 
    }
}

function toggleNav() {
    if (mainSideBar.style.width === "250px") {
        closeNav();
    } else {
        openNav();
    }
}






