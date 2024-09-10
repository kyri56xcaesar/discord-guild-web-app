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

client = discord.Client(intents=intents)

my_guilds = ['$DADS', 'ΗΜΜΥ']
ROLE_INDEX = 'HOF'


@client.event
async def on_ready():
    print(f'Logged in as {client.user}')
    
    data = await gatherData(verbose=False)
         
    url = 'http://localhost:6969/guild/members'
    headers = {
        'user-agent':'my_discord_app/0.0.1',
        'content-type':'application/json'
    }
    response = requests.post(url=url, headers=headers, data=json.dumps(data))
    
    print(response.status_code, response.text)
        


async def gatherData(verbose=True):
    
    now = datetime.now()
    after = now.replace(month=now.month - 2)
    before = now.replace(month=now.month - 1)
    
 
        # Loop through all servers (guilds) the bot is part of
    for guild in client.guilds:
        #print(f"\nServer Name: {guild.name}")
        
        
            
        
        # Perform calculations and send data to my service
        # Need to count messages per member for all channels and sort them
        # send data to service as in member, name, avatar, avatar_decor, banner, color, joined_at
        data = list()
        for member in guild.members:
            
            if ROLE_INDEX in [role.name for role in member.roles]:
            
                m_data = {
                    'user':member.name,
                    'nick': member.nick if member.nick is not None else 'None',
                    'avatar':str(member.avatar),
                    'joined_at':member.joined_at.isoformat(),  
                    'status':member.status.value,
                    'roles':[role.name for role in member.roles],
                    'msg_count':0
                }
                for channel in guild.text_channels:
                    
                    async for message in channel.history(before=before, after=after):
                        if message.author == member:
                            m_data['msg_count'] += 1
                            if verbose:
                                print(f'Member: {member.name} -> {channel} -> message: {message.created_at}\t{message.content}')

                if verbose:
                    print(m_data)
                data.append(m_data)
        
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