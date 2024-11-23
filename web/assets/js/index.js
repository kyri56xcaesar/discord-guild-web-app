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

// Utility functions
function formatDate(dateString) {
    const date = new Date(dateString);

    return date.toLocaleDateString('en-US', {year: 'numeric', month: 'long', day: 'numeric'})
    
}

function hexToRGBA(hex_string, alpha = 1) {

    isHexColor = hex => typeof hex === 'string' && hex.length === 6 && !isNaN(Number('0x' + hex))

    if (!isHexColor(hex_string)) {
        return hex_string;
    }

    hex_string = hex_string.replace(/^#/, '');

    if (hex_string.length === 3) {
        hex = hex.split('').map(char => char + char).join('');
    }

    const r = parseInt(hex_string.substring(0, 2), 16);
    const g = parseInt(hex_string.substring(2, 4), 16);
    const b = parseInt(hex_string.substring(4, 6), 16);

    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
}

