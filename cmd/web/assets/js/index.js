const sideBarAffectedElements = document.querySelectorAll('.sidebar-affected');
// console.log(sideBarAffectedElements)

const mainSideBar = document.getElementById("main-siderbar");
const navbtn = document.getElementById("navbtn");
const togglebtn = document.getElementById("togglebtn")

function openNav() {
    mainSideBar.style.width = "250px";
    navbtn.style.marginLeft = "20px";
    togglebtn.innerText = "Less";

    for (let element of sideBarAffectedElements) {
        element.style.marginLeft = '250px'; 
    }
}

function closeNav() {
    mainSideBar.style.width = "0";
    navbtn.style.marginLeft = "0";
    togglebtn.innerText = "More";

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


const botsDisplay = document.getElementById("bots-display");
const servicesDisplay = document.getElementById("hof-display");
const clientsDisplay = document.getElementById("clients-display");

// Main Body display options
function showClients() {

    // console.log('Displaying Clients content...')

    botsDisplay.style.display = "none";
    servicesDisplay.style.display = "none"
    clientsDisplay.style.display = "block"
}

function showHof() {

    // console.log('Displaying Services content...')
    
    botsDisplay.style.display = "none";
    servicesDisplay.style.display = "block"
    clientsDisplay.style.display = "none"
}

function showMain() {
    
    // console.log('Displaying main content...')

    botsDisplay.style.display = "block";
    servicesDisplay.style.display = "none"
    clientsDisplay.style.display = "none"
}






