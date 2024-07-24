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

// Footer redirection
document.getElementById('contact-ref').addEventListener('click', function(event) {
    event.preventDefault();
    closeNav()
    window.scrollTo({ top: document.body.scrollHeight, behavior: 'smooth' });
});






// img rotation
document.addEventListener("DOMContentLoaded", function() {
    const rotatableImage = document.querySelector('.rotatable-image');
    let isDragging = false;
    let previousAngle = 0;
    let startX, startY;

    function getAngle(x, y) {
        return Math.atan2(y, x) * (180 / Math.PI);
    }

    rotatableImage.addEventListener('mousedown', function(event) {
        isDragging = true;
        const rect = rotatableImage.getBoundingClientRect();
        startX = event.clientX - rect.left - rect.width / 2;
        startY = event.clientY - rect.top - rect.height / 2;
        previousAngle = getAngle(startX, startY);
        rotatableImage.style.cursor = 'grabbing';
    });

    document.addEventListener('mousemove', function(event) {
        if (isDragging) {
            const rect = rotatableImage.getBoundingClientRect();
            const currentX = event.clientX - rect.left - rect.width / 2;
            const currentY = event.clientY - rect.top - rect.height / 2;
            const currentAngle = getAngle(currentX, currentY);
            const angleDifference = currentAngle - previousAngle;

            rotatableImage.style.transform = `rotate(${angleDifference}deg)`;
        }
    });

    document.addEventListener('mouseup', function() {
        if (isDragging) {
            isDragging = false;
            rotatableImage.style.cursor = 'grab';
        }
    });

    rotatableImage.addEventListener('mouseleave', function() {
        if (isDragging) {
            isDragging = false;
            rotatableImage.style.cursor = 'grab';
        }
    });
});



