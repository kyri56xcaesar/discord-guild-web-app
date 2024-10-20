import discord
import json
import requests
import asyncio
import logging
from dotenv import dotenv_values
from discord.ext import commands, tasks
from datetime import datetime, timedelta

config = dotenv_values("../../bots.env")
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s', 
                    handlers=[logging.FileHandler("bot.log"), logging.StreamHandler()])
logger = logging.getLogger(__name__)

# Set up intents


# Create bot instance with command prefix '!' and intents
bot = commands.Bot(command_prefix="!", intents=discord.Intents.all())

ROLE_INDEX = config['ROLE_INDEX']
URL = 'http://' + config['IP'] + ':' + config['PORT'] + '/guild/members'




# UTILS
def int_to_hex_color(integer_value):
    if not (0 <= integer_value <= 0xFFFFFF):
        raise ValueError("The integer value must be between 0 and 16777215 (0xFFFFFF).")
    
    hex_value = hex(integer_value)[2:]
    hex_value = hex_value.zfill(6)
    
    return f'#{hex_value}'




# # # EVENTS # # #
@bot.event
async def on_ready():
    logger.info(f'Logged in as {bot.user}')
    # Start the server in a background task
    await bot.loop.create_task(server_thread())



# # # COMMANDS # # # 
# Bot command to manually trigger data gathering and forwarding
@bot.command(name="gather")
async def gather_command(ctx, monthOffset: int = 1):
    try:
        logger.info(f"Gather command invoked with monthOffset: {monthOffset}")
        data = await gatherData(verbose=True, monthOffset=monthOffset)
        await ctx.send(f"Gathered data for {monthOffset} months.")
        # Optionally forward data here
        response = await forwardData(data, URL)
        await ctx.send(f"Data forwarded. Response: {response.status_code}")
    except Exception as e:
        logger.error(f"Error in gather_command: {e}", exc_info=True)
        await ctx.send("An error occurred while gathering data.")

# If you want the bot to have a periodic task, you can set up a loop:
@tasks.loop(hours=24)
async def periodic_data_collection():
    data = await gatherData(verbose=False, monthOffset=1)
    response = await forwardData(data, URL)
    print(f"Data collected and forwarded. Status: {response.status_code}")








# Functions
async def forwardData(data, url):
    headers = {
        'user-agent': 'my_discord_app/0.0.1',
        'content-type': 'application/json; charset=utf-8'
    }
    response = requests.post(url=url, headers=headers, data=json.dumps(data, ensure_ascii=False).encode('utf-8'))
    return response


async def gatherData(verbose=True, monthOffset=-1):
    before = datetime.now()

    if monthOffset > 0:
        after = before - timedelta(days=monthOffset * 30)
    else:
        after = None

    data = dict()

    # Loop through all servers (guilds) the bot is part of
    for guild in bot.guilds:
        member_data = list()
        # Find HOF members
        for member in guild.members:
            if ROLE_INDEX in [role.name for role in member.roles]:
                m_data = {
                    'userguild': guild.name,
                    'username': member.name,
                    'nickname': member.nick if member.nick is not None else 'None',
                    'avatarurl': str(member.avatar),
                    'displayavatarurl': str(member.display_avatar),
                    'bannerurl': member.banner.url if member.banner is not None else 'None',
                    'usercolor': str(int_to_hex_color(member.color.value)).upper(),
                    'joinedat': str(member.joined_at),
                    'userstatus': member.raw_status,
                    'userroles': [{'role_name': role.name, 'role_color': str(int_to_hex_color(role.color.value).upper())} for role in member.roles],
                    'usermessages': list(),
                    'messagecount': 0
                }
                member_data.append(m_data)

            if member.bot:
                bot_data = {
                    'bot': True,
                    'userguild': guild.name,
                    'username': member.name,
                    'avatarurl': str(member.avatar),
                    'joinedat': str(member.joined_at),
                    'userstatus': member.raw_status,
                    'userroles': [{'role_name': role.name, 'role_color': str(int_to_hex_color(role.color.value).upper())} for role in member.roles],
                    'usermessages': list(),
                    'messagecount': 0
                }
                member_data.append(bot_data)

        data[guild.name] = member_data

        # Retrieve 1 month message history
        for channel in guild.text_channels:
            async for message in channel.history(before=before, after=after, limit=None):
                for item in data[guild.name]:
                    if item['username'] == message.author.name:
                        item['messagecount'] += 1
                        item['usermessages'].append({
                            "content": message.clean_content,
                            "channel": message.channel.name,
                            "createdat": message.created_at
                        })
                        break  # Stop after finding the correct user

    return data







# # # SERVER # # #
async def handle_client(reader, writer):
    logger.info("Client connected")
    try: 
        data = await reader.read(1024)
        signal = data.decode('utf-8')

        logger.info(f"Received signal: {signal}")

        if signal == 'gather':
            data = await gatherData(verbose=False, monthOffset=1)
            logger.info(data)
            #response_data = json.dumps(data).encode('utf-8')
            #logger.info(response_data)
            logger.info("Data sent back to client.")

        elif signal == 'check':
            writer.write("I am here!".encode('utf-8'))
            await writer.drain()
            logger.info("Sent 'I am here!' response.")

    except Exception as e:
        logger.error(f"An error occurred: {e}", exc_info=True)
    finally:
        writer.close()
        await writer.wait_closed()
        logger.info("Client connection closed.")


async def server_thread():
    lport = int(config['LPORT'])  # Convert port to integer

    try:
        server = await asyncio.start_server(handle_client, 'localhost', lport)
        addr = server.sockets[0].getsockname()
        logger.info(f'Serving on {addr}')

        async with server:
            await server.serve_forever()
    except Exception as e:
        logger.error(f"Failed to start the server: {e}", exc_info=True)





if __name__ == "__main__":
    token = config['snake_bot_token']
    bot.run(token)






