
setHofTitle(guild_name);
function setHofTitle(guild_name) {

    const matrix_switch = document.getElementById('matrix-switch');
    const title_effect_switch = document.getElementById('title-effect-switch');

    matrix_switch.setAttribute('checked', 'checked');
    title_effect_switch.setAttribute('checked', 'checked');


    let hofTitle = document.getElementById('hof-title');
    let currentDate = new Date();

    currentDate.setMonth(currentDate.getMonth() - 1);
    let month = currentDate.getMonth() + 1;
    let year = currentDate.getFullYear();
    hofTitle.innerText = `${guild_name} HALL OF FAME - ${month}/${year}`;
}


function toggleTitleAnimation() {

    // console.log('Tittle effect toggle');
    const title = document.getElementById('hof-title');

    if (title.style.animationPlayState === 'paused') {
        title.style.animationPlayState = 'running'; // Resume the animation
    } else {
        title.style.animationPlayState = 'paused';  // Pause the animation (freeze)
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


