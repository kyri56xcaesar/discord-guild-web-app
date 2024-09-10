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
      console.log(data);
      displayMembers(data);
      
}

function displayMembers(data) {
    const hof_list = document.getElementById('hof-list');

    data.forEach(member => {
        let li = document.createElement('li');

        li.innerHTML = `
            <div class="hof-entry">
                <img src="${member.avatar}" alt="${member.member}'s avatar" width="150px" />
                <h3>${member.user} (${member.nick || 'No Nick'})</h3>
                <p>Message Count: ${member.msg_count}</p>
                <p>Status: ${member.status}</p>
                
                <p>Roles: ${member.roles.join(", ")}</p>
                <p>Joined At: ${member.joined_at}</p>
            </div>
            <hr>
        `;

        hof_list.appendChild(li);
    })
}

fetchMembers();
