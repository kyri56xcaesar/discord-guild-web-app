var guild_name = '$DADS';

let headersList = {
    "Accept": "*/*",
    "Content-Type": "application/json"
   }

async function fetchMembers() {
    let response = await fetch("http://localhost:6969/guild/members", { 
        method: "GET",
        headers: headersList
      });
      
      let data = await response.json();

      guild_name = data[0].guild;
      data.sort(function(first, second) {
        return second.msg_count - first.msg_count
      })

      console.log(data);
      displayMembers(data);
      
}

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



function formatDate(dateString) {
    const date = new Date(dateString);

    return date.toLocaleDateString('en-US', {year: 'numeric', month: 'long', day: 'numeric'})
    
}

function displayMembers(data) {
    const hof_list = document.getElementById('hof-list');

    data.forEach((member, index) => {
        let li = document.createElement('li');
        

        li.innerHTML = `
            <div class="hof-entry-left">
                <h1>${index + 1}.</h1>
                <img class="border-color-${member.status}"src="${member.avatar}" alt="${member.member}'s avatar" />
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

fetchMembers();


// Graphics buttons
function toggleMatrixEffect(){
    console.log();
}

function toggleTitleEffect() {
    console.log()
}

