import discord
from datetime import datetime, timedelta
from dotenv import dotenv_values

client = discord.Client(intents=discord.Intents.default())


config = dotenv_values(".env")

token = config["stratigos_bot_token"]
vip_list = config["vip_list_user_ids"].split(",")

print(f"token: {token}, vip_list: {vip_list}")


phrase = "\nH wra poy to voulwneis"
responding = True
HELP = "GENERALLLLLLL\n\
        -voulwne : to voulwnw\n\
        -peripolia : prosexws\n\
        -vohtheia : auti i malakia p molis eides.\n"


@client.event
async def on_ready():
    print("We have logged as {0.user}".format(client))




@client.event
async def on_message(message):

    global responding

    if message.author == client.user:
        return

    if message.author.bot:
        return

    if message.author.id in vip_list:
        return


    msg = message.content

    if msg.startswith("-voulwne"):

        responding = False
        await message.channel.send("ANAPAUSI")
        return

    if msg.startswith("-peripolia"):

        responding = True
        await message.channel.send("PROSOXI")
        return


    if msg.startswith("-vohtheia"):
        await message.channel.send(HELP)
        return

    now = datetime.now() + timedelta(hours=0)
    ts = now.strftime("%H:%M")


    if responding:
        await message.channel.send(f"{ts} {message.author.mention} {phrase}")
    

client.run(token)




