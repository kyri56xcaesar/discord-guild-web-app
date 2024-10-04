import discord
import json
import requests

from dotenv import dotenv_values
from datetime import datetime, timedelta


config = dotenv_values("../../bots.env")


intents = discord.Intents.default()
intents.members = True
intents.guild_messages = True
intents.message_content = True
intents.messages = True
intents.presences = True
intents.guilds = True


client = discord.Client(intents=intents)

ROLE_INDEX = config['ROLE_INDEX']
URL = 'http://'+config['IP']+':'+config['PORT']+'/guild/members'


def int_to_hex_color(integer_value):
    # Ensure the integer is within the valid range for colors (0x000000 to 0xFFFFFF)
    if not (0 <= integer_value <= 0xFFFFFF):
        raise ValueError("The integer value must be between 0 and 16777215 (0xFFFFFF).")
    
    hex_value = hex(integer_value)[2:]
    
    
    hex_value = hex_value.zfill(6)
    
    return f'#{hex_value}'


@client.event
async def on_ready():
    print(f'Logged in as {client.user}')
    
    # service url
    
    # perhaps make it multithreaded
    data = await gatherData(verbose=False, monthOffset=6)
    with open('test.json', 'w', encoding='utf-8') as f:
        json.dump(fp=f, obj=data, ensure_ascii=False)
    # multithreaded as well
    
    response = await forwardData(data['$Dads'], URL)
    # print(response.status_code, response.text)

    
async def forwardData(data, url):
    headers = {
        'user-agent':'my_discord_app/0.0.1',
        'content-type':'application/json; charset=utf-8'
    }
    response = requests.post(url=url, headers=headers, data=json.dumps(data, ensure_ascii=False).encode('utf-8'))
    
    
    return response
        

async def gatherData(verbose=True, monthOffset=-1):
    
    before = datetime.now()
    
    if monthOffset > 0:
        after = before - timedelta(days=monthOffset*30)
    else: 
        after = None
    
    data = dict()
    
    # Loop through all servers (guilds) the bot is part of
    for guild in client.guilds:
        member_data = list()
        # Find HOF members
        for member in guild.members:
            
            if ROLE_INDEX in [role.name for role in member.roles]:
            
                m_data = {
                    'guild':guild.name,
                    'user':member.name,
                    'nick': member.nick if member.nick is not None else 'None',
                    'avatar':str(member.avatar),
                    'display_avatar':str(member.display_avatar),
                    'banner':member.banner.url if member.banner is not None else 'None',
                    'user_color':str(int_to_hex_color(member.color.value)).upper(),
                    'joined_at':member.joined_at.isoformat(),
                    'status':member.raw_status,
                    'roles':[{'role_name': role.name, 'role_color':str(int_to_hex_color(role.color.value).upper())} for role in member.roles],
                    'messages':list(),
                    'msg_count':0
                }
                member_data.append(m_data)
                
                
            if member.bot:
                
                bot_data = {
                    'bot':True,
                    'guild':guild.name,
                    'user':member.name,
                    'avatar':str(member.avatar),
                    'joined_at':member.joined_at.isoformat(),
                    'status':member.raw_status,
                    'roles':[{'role_name': role.name, 'role_color':str(int_to_hex_color(role.color.value).upper())} for role in member.roles],
                    'messages':list(),
                    'msg_count':0
                }
                member_data.append(bot_data)

        data[guild.name] = member_data
        
        # Retrieve 1 month message history
        for channel in guild.text_channels:           
            async for message in channel.history(before=before, after=after, limit=None):
                # print(f'Member: {member.name} -> {channel} -> message: {message.created_at}\t{message.content}')
                for item in data[guild.name]:
                    if item['user'] == message.author.name:
                    # Increment the 'messages' count by 1
                        item['msg_count'] += 1
                        item['messages'].append(f'{message.clean_content}, {message.created_at}')
                        break  # Stop after finding the correct user

        
        
        return data
                       

if __name__ == "__main__":
    #app.run(host='0.0.0.0', port=PORT, debug=True)
    

    token = config['snake_bot_token']
    # print(f'TOKEN: {token}')
    client.run(token)

