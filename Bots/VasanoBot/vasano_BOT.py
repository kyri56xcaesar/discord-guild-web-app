import discord
import random
from dotenv import load_dotenv
import os, sys

parent_directory = os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir))
sys.path.append(parent_directory)

import bot_utils as ut

load_dotenv()

token = os.getenv("bot_token")
userID = os.getenv("userID")
general_ch_id = os.getenv("general_ch_id")
bot_testing_ch_id = os.getenv("bot_testing_ch_id")


print(f"token: {token}, userID: {userID}, general_channel: {general_ch_id}")


intents = discord.Intents.all()
client = discord.Client(intents=intents)

activated = False
trigs = ["Kalhspera kai kali vradia edw me ton aderfo fofota", "gamo ton karioli ton OTE apo tis 7 to prwi mexri twra den eixa internet gamo tinpanagia tous prwi prwi",
         "gamw tis manes tous", "simera anastithike o bines", "me trexoyn oi traktores", "ksekinise to praktorio", "ff x9 matchfixers/griefers/hostagers", "epesa PLAT apo master 230 lp se mia nuxta", "pame vraxous gt me trexoyn ta mpastarda", "koke gamw tinmana sou vlaka",
         "me treksan tapromo mexri diamond 0 lp", "legit ama den paw master to kleinw kai 3 meres naparei ", "na se paw mia volta aptis kuries?", "o mikros apo tin pisw porta", "me trexei o lazer"]

when = ["yorick", "praktoras"]


tracker = list()

counter = 0


# UTILITY FUNCS

def print_pops():
    message = ""

    message += "##################\n"
    for i, pop in enumerate(trigs):

        message += f"{i+1}--> {pop}.\n"
    message += "##################\n"

    return message


@client.event
async def on_ready():
    print("Logged in as ")
    print(client.user.name)
    print(client.user.id)
    print("--------")

    channel = await client.fetch_channel(general_ch_id)

    print(channel)

    await channel.send(content="Kalhspera kai kali vradia edw me ton aderfo fofota", delete_after=10)


@client.event
async def on_message(msg):

    # Check if msg is fetched
    if msg is None:
        return
    # Check if message was sent by me.
    if msg.author == client.user:
        return

    # Get string content.
    message = msg.content

    # Access global variable
    # To respond or not
    global activated, counter

    if activated:

        # counter += 1
        # if counter == 3:
        #     await msg.channel.send(random.choice(trigs))
        #     counter = 0
        if message in when:
            await msg.channel.send(random.choice(trigs))

    # Control messages
    if message.startswith("%vale"):
        m = message.split("%vale", maxsplit=1)[1]
        trigs.append(m)
    elif message.startswith("%vgale"):
        m = message.split("%vgale", maxsplit=1)[1]
        if m.isnumeric():
            index = int(m)
            if index <= len(trigs) and index <= 1:
                del trigs[index - 1]
            else:
                # send invalid index warning
                pass
        elif m in trigs:
            trigs.remove(m)
        else:
            # send 'invalid' warning
            pass
    elif message.startswith("ti exoyme edw"):
        await msg.channel.send(print_pops())
    elif message.startswith("%voulwne"):
        if activated == False:
            await msg.channel.send("eimai seri")
        else:
            activated = False
            await msg.channel.send("paw gia upno")

    elif message.startswith("%ksupna"):
        if activated:
            await msg.channel.send("Edw eimai re clown")
        else:
            activated = True
            await msg.channel.send("Ksekinaw to treksimo.")

    elif message.startswith("%vohtheia"):
        await msg.channel.send("TI VOHTHEIA MWRI KURIA")

    if msg.author.id == id:
        tracker.append(message)


client.run(token)


ut.save_data("Vassano", tracker)
