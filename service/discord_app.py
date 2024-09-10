import os
import discord
import json
import requests
from dotenv import load_dotenv
from datetime import datetime


load_dotenv()

intents = discord.Intents.default()
intents.members = True
intents.guild_messages = True
intents.message_content = True
intents.messages = True
intents.presences = True

client = discord.Client(intents=intents)

my_guilds = ['$DADS', 'ΗΜΜΥ']
ROLE_INDEX = 'HOF'
URL = 'http://localhost:6969/guild/members'


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
    data = await gatherData(verbose=False, monthOffset=1)
    with open('test.json', 'w') as f:
        json.dump(fp=f, obj=data)
    # multithreaded as well
    
    response = await forwardData(data, URL)
    print(response.status_code, response.text)

    
async def forwardData(data, url):
    headers = {
        'user-agent':'my_discord_app/0.0.1',
        'content-type':'application/json'
    }
    response = requests.post(url=url, headers=headers, data=json.dumps(data))
    
    
    return response
        

async def gatherData(verbose=True, monthOffset=-1):
    
    now = datetime.now()
    after = now.replace(month=now.month - (monthOffset+1), day=1)
    
    if monthOffset >= 0:
        before = now.replace(month=now.month - monthOffset, day=1)
    else:
        before = now
    
    # Loop through all servers (guilds) the bot is part of
    for guild in client.guilds:
        data = list()

        # Find HOF members
        for member in guild.members:
            
            if ROLE_INDEX in [role.name for role in member.roles]:
            
                m_data = {
                    'guild':'$DADS',
                    'user':member.name,
                    'nick': member.nick if member.nick is not None else 'None',
                    'avatar':str(member.avatar),
                    'banner':member.banner.url if member.banner is not None else 'None',
                    'user_color':str(int_to_hex_color(member.color.value)).upper(),
                    'joined_at':member.joined_at.isoformat(),
                    'status':member.raw_status,
                    'roles':[{'role_name': role.name, 'role_color':str(int_to_hex_color(role.color.value).upper())} for role in member.roles],
                    'msg_count':0
                }
                data.append(m_data)

        
        # Retrieve 1 month message history
        for channel in guild.text_channels:           
            async for message in channel.history(before=before, after=after, limit=1000):
                # print(f'Member: {member.name} -> {channel} -> message: {message.created_at}\t{message.content}')
                for item in data:
                    if item['user'] == message.author.name:
                    # Increment the 'messages' count by 1
                        item['msg_count'] += 1
                        break  # Stop after finding the correct user

        
        
        return data
                       

if __name__ == "__main__":
    #app.run(host='0.0.0.0', port=PORT, debug=True)
    

    token = os.getenv('DISCORD_TOKEN')
    print(f'TOKEN: {token}')
    client.run(token)



#   # test_data
#     data = [
#         {
#             'id':0, 
#             'user': 'kitsunee0', 
#             'nick': 'None', 
#             'avatar': 'https://cdn.discordapp.com/avatars/487049871950610447/a_c7ba5592e738cc81f73b5a525a99d347.gif?size=1024', 
#             'joined_at': '2020-07-22T20:11:05.319000+00:00', 
#             'status': 'offline', 
#             'roles': [
#                 '@everyone', 
#                 'HOF', 
#                 'PR', 
#                 'Ο ΓΑΜΙΑΣ ΤΗΣ ΓΕΙΤΟΝΙΑΣ', 
#                 'IronImport', 
#                 'KEKW', 
#                 'Lepra sta daxtula', 
#                 'KOURADOPAIKTIS', 
#                 'Fortnite', 
#                 'VaLoRaNt pRo', 
#                 'imPOSTER', 
#                 'Ελο Σλέιβ', 
#                 '--------------GAMER---------------', 
#                 'ζογκλέρ', 
#                 '---------------ROLES-----------------', 
#                 'BLACK BULL E-SPORTS', 
#                 'Candy Giver', 
#                 'VISMA PASS', 
#                 'still Σκουπίδια', 
#                 'CLOWN + PIG combo', 
#                 'GENE DIFF VICTIMS', 
#                 'pissrat pissrandom pissdogs', 
#                 'mini commadakia', '(-)S1GMA D0GOS'
#                 ], 
#             'msg_count': 0}, 
#         {
#             'id':0, 
#             'user': 'rank1to', 
#             'nick': 'codfish client', 
#             'avatar': 'https://cdn.discordapp.com/avatars/429980662355722240/8545e0173186b672ad38662b9473901c.png?size=1024', 
#             'joined_at': '2023-08-11T20:47:18.231000+00:00', 
#             'status': 'offline', 
#             'roles': [
#                 '@everyone', 
#                 'HOF', 
#                 'TILTED', 
#                 'Σεριστας', 
#                 'BLACKBULLSGAMIADES', 
#                 '"Master elo player "', 
#                 'VaLoRaNt pRo', 
#                 'imPOSTER', 
#                 'Ελο Σλέιβ', 
#                 '🛡️ VALHEIM 🛡️', 
#                 'Σαπόρτ', 
#                 '🃏JackOfAllTrades🃏', 
#                 'BLACK BULL E-SPORTS', 
#                 '--------------TEAM-----------------', 
#                 'THIRD IMPOSTOR', 
#                 'OG HADES', 
#                 'THRESH-ABLE', 
#                 'Bigg Choker', 
#                 'L FUCKING 9', 
#                 '--------------Attribute--------------', 
#                 'Θείος στην δίωξη', 
#                 'MATRIX NPC RANDOMS', 
#                 'ANDRIKO MORIO ENJOYERS', 
#                 'ΔELTA COURIER BOSS', 
#                 'Πατεράδες', 
#                 'ΩμEGA CHIEF BOSS'
#                 ], 
#             'msg_count': 11}
#     ]