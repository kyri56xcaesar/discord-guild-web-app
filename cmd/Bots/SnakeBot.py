import discord
import json
import requests
import asyncio
import logging
from dotenv import dotenv_values
from discord import app_commands
from discord.ext import commands, tasks
from datetime import datetime, timedelta

config = dotenv_values("../../bots.env")
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s', 
                    handlers=[logging.FileHandler("bot.log"), logging.StreamHandler()])
logger = logging.getLogger(__name__)


# constants
MAIN_GUILD_ID=config['$dads_guild_id']
MC_DISPLAY_ID=config['member_count_channel_id']

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
    try:
        synced = await bot.tree.sync()
        logger.info(f"Synced {len(synced)} command(s)")
        await update_member_count(guild_id=MAIN_GUILD_ID, channel_id=MC_DISPLAY_ID)
    except Exception as e:
        logger.error(f"Failed to sync {e}")
    await bot.loop.create_task(server_thread())
    
    
async def update_member_count(guild_id, channel_id):
    
    guild = bot.get_guild(guild_id)
    if not guild:
        return
    
    member_count = guild.member_count
        
    logging.info(f"Member count is: {member_count}")
    channel = bot.get_channel(channel_id)
    logging.info(f"Channel to be changed is: {channel.name}")

    if channel is None:
        # If the channel is not found in cache, fetch it
        try:
            channel = await bot.fetch_channel(channel_id)
        except discord.NotFound:
            logger.info(f"Channel with ID {channel_id} not found.")
            return
        except discord.Forbidden:
            logger.info(f"Insufficient permissions to access channel ID {channel_id}.")
            return
        except Exception as e:
            logger.info(f"An error occurred while fetching the channel: {e}")
            return
    
    
    await channel.edit(name=f"Humans: {member_count}")


# # # COMMANDS # # # 
# Bot command to manually trigger data gathering and forwarding
@bot.tree.command(name="gather")
@app_commands.describe(months="How many months?")
async def gather_command(interaction: discord.Interaction, months: str):
    try:
        await interaction.response.send_message(f"{interaction.user.name} said: `{months}`")
        logger.info(f"Gather command invoked with monthOffset: {months}")
        data = await gatherData(verbose=True, monthOffset=int(months))
        # logger.info(f"Data gathered:  {data }")
        logger.info(f"Gathered data for {months} months.")
        # Optionally forward data here
        response = await forwardData(data, URL)
        logger.info(f"Data forwarded. Response: {response.status_code}")
    except Exception as e:
        logger.error(f"Error in gather_command: {e}", exc_info=True)
        logger.error("An error occurred while gathering data.")

# If you want the bot to have a periodic task, you can set up a loop:
# @tasks.loop(hours=24)
# async def periodic_data_collection():
#     data = await gatherData(verbose=False, monthOffset=1)
#     response = await forwardData(data, URL)
#     logger.info(f"Data collected and forwarded. Status: {response.status_code}")







# Functions
async def forwardData(data, url):
    headers = {
        'user-agent': 'my_discord_app/0.0.1',
        'content-type': 'application/json; charset=utf-8'
    }
    response = requests.post(url=url, headers=headers, data=json.dumps(data, ensure_ascii=False, indent=4, sort_keys=True, default=str).encode('utf-8'))
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
                    'joinedat': str(member.joined_at.strftime('%Y-%m-%d')),
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
                    'joinedat': str(member.joined_at.strftime('%Y-%m-%d')),
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
        while True:
            data = await reader.readline()
            if not data:
                # client here
                break
            if len(data) > 2048:
                response = "Error: Message too long."
                writer.write((response + '\n').encode('utf-8'))
                await writer.drain()
                continue  
            signal = data.decode('utf-8').strip()
            logger.info(f"Received signal: {signal}")
            
            channel_id = 778642179337093160  # Replace with the actual channel ID
            channel = bot.get_channel(channel_id)

            if channel is None:
                # If the channel is not found in cache, fetch it
                try:
                    channel = await bot.fetch_channel(channel_id)
                except discord.NotFound:
                    logger.error(f"Channel with ID {channel_id} not found.")
                    continue  # Skip sending the message
                except discord.Forbidden:
                    logger.error(f"Insufficient permissions to access channel ID {channel_id}.")
                    continue
                except Exception as e:
                    logger.error(f"An error occurred while fetching the channel: {e}")
                    continue

            try:
                # Send the message to the channel
                await channel.send(f"{signal}")
                logger.info(f"Sent message to channel ID {channel_id}.")
            except discord.Forbidden:
                logger.error(f"Insufficient permissions to send messages to channel ID {channel_id}.")
            except Exception as e:
                logger.error(f"An error occurred while sending message to the channel: {e}")


            if signal == 'gather':
                data = await gatherData(verbose=False, monthOffset=1)
                response_data = json.dumps(data)
                writer.write((response_data + '\n').encode('utf-8'))
                await writer.drain()
                logger.info("Data sent back to client.")

            elif signal == 'check':
                response = "I am here!"
                writer.write((response + '\n').encode('utf-8'))
                await writer.drain()
                logger.info("Sent 'I am here!' response.")

            elif signal == 'exit':
                logger.info("Client requested to close the connection.")
                break  

            else:
                response = "Unknown command."
                writer.write((response + '\n').encode('utf-8'))
                await writer.drain()
                logger.info("Sent 'Unknown command' response.")
                
            

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

