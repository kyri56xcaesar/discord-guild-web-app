import os
import discord
from dotenv import load_dotenv



load_dotenv()

intents = discord.Intents.default()
intents.members = True

client = discord.Client(intents=intents)

my_guilds = ['$DADS', 'ΗΜΜΥ']


@client.event
async def on_ready():
    print(f'Logged in as {client.user}')
    
    # Loop through all servers (guilds) the bot is part of
    for guild in client.guilds:
        print(f"\nServer Name: {guild.name}")
        
        
        # Loop through all channels in the server
        print('Channels:')
        for channel in guild.channels:
            print(f' - {channel.name} ({channel.type})')
            
        # Loop through all members in the server
        print('Members:')
        for member in guild.members:
            print(f' - {member.name}#{member.discriminator}')




if __name__ == "__main__":
    #app.run(host='0.0.0.0', port=PORT, debug=True)
    

    token = os.getenv('DISCORD_TOKEN')
    print(f'TOKEN: {token}')
    client.run(token)
