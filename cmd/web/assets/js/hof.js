var guild_name = '$DADS';
IP = 'localhost';
// IP = '85.75.213.54'
PORT = '6969'


// Utility functions
function formatDate(dateString) {
    const date = new Date(dateString);

    return date.toLocaleDateString('en-US', {year: 'numeric', month: 'long', day: 'numeric'})
    
}

function hexToRGBA(hex_string, alpha = 1) {
    hex_string = hex_string.replace(/^#/, '');

    if (hex_string.length === 3) {
        hex = hex.split('').map(char => char + char).join('');
    }

    const r = parseInt(hex_string.substring(0, 2), 16);
    const g = parseInt(hex_string.substring(2, 4), 16);
    const b = parseInt(hex_string.substring(4, 6), 16);

    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
}



// Include axios if you're using a CDN
// <script src="https://cdn.jsdelivr.net/npm/axios/dist/axios.min.js"></script>

async function fetchMembers() {
    try {
        const response = await axios.get("http://"+IP+":"+PORT+"/guild/members", {
            headers: {
                "Content-Type": "application/json"
            }
        });
        
        data = response.data;
        guild_name = data[0].guild;
        data.sort(function(first, second) {
            return second.msg_count - first.msg_count
        })

        console.log(data);
        displayMembers(data);

    } catch (error) {
        // Log detailed error information
        if (error.response) {
            console.error("Error in response:", error.response.status, error.response.data);
        } else if (error.request) {
            console.error("Error in request:", error.request);
        } else {
            console.error("General Error:", error.message);
        }
    }
}

fetchMembers();






// Display functions and setups
function setMainTitle(guild_name) {
    let mainTitle = document.getElementById('main-title');

    mainTitle.innerText = `${guild_name}`
}

setMainTitle(guild_name);

function setHofTitle(guild_name) {
    let hofTitle = document.getElementById('hof-title');
    let currentDate = new Date();

    currentDate.setMonth(currentDate.getMonth() - 1);
    let month = currentDate.getMonth() + 1;
    let year = currentDate.getFullYear();
    hofTitle.innerText = `${guild_name} HALL OF FAME - ${month}/${year}`;
}

setHofTitle(guild_name);




function displayMembers(data) {
    const hof_list = document.getElementById('hof-list');

    data.forEach((member, index) => {
        let li = document.createElement('li');
        li.setAttribute('id', `hof-entry-${index + 1}`);
        li.setAttribute('class', 'hof-entry')

        li.style=`background-color: ${hexToRGBA(member.user_color, 0.08)}; background: url(${member.banner})`

        // console.log('avatar: ' + member.avatar + 'banner: ' + member.banner);

        if (index != 0 && index != 1 && index != 2) li.style.borderBottom = '1px solid #ccc';

        const default_icon = "assets/mozart_draw.jpg";

        li.innerHTML = `
            <div class="hof-entry-left">
                <h1>${index + 1}.</h1>
                <img class="border-color-${member.status}"src="${member.avatar === 'None' ? default_icon : member.avatar}" alt="${default_icon}" />
                <h3>${member.user} <br>(${member.nick || 'No Nick'})</h3>
            </div>

            <div style="border-left:1px solid #000;height:250px"></div>
            
            <div class="hof-entry-right">
                <p><strong>Message Count:</strong> <span class="msg-count">${member.msg_count}</span></p>
                
                <strong>Roles:</strong>
                <p class="roles-display" id="roles-display"></p>
                
                <p><strong>Joined At:</strong> ${formatDate(member.joined_at)}</p>
                
                <p>
                    <strong>Status:</strong> ${member.status}
                    <span class="status-light ${member.status}"></span>
                </p>
            </div>
            <hr>
        `;
        hof_list.appendChild(li);
        
        let rolesDisplay = li.querySelector('.roles-display');

        let rolesHtml = member.roles.map(role => {
            return `<span class="roles-full" style="color: ${role.role_color};">${role.role_name}</span>`;
        }).join(", ");
  
        rolesDisplay.innerHTML = rolesHtml;
 
        

    })

    
}

const title = document.getElementById('hof-title');


function toggleTitleAnimation() {

    // console.log('Tittle effect toggle');
    
    if (title.style.animationPlayState === 'paused') {
        title.style.animationPlayState = 'running'; // Resume the animation
    } else {
        title.style.animationPlayState = 'paused';  // Pause the animation (freeze)
    }
}


const matrix_switch = document.getElementById('matrix-switch');
const title_effect_switch = document.getElementById('title-effect-switch');

matrix_switch.setAttribute('checked', 'checked');
title_effect_switch.setAttribute('checked', 'checked');

