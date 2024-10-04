import discord
import requests
import json
import random
import bot_utils as ut

from dotenv import dotenv_values




intents = discord.Intents.all()
client = discord.Client(intents=intents)


config = dotenv_values("../../bots.env")

token = config["koke_bot_token"]
id = config["koke_user_id"]
tracker = list()

print(f"token: {token}, userID: {id}")


trigger_words = ["Koke", "koke", "KOKE", "ela", "0/100", "3 cs", "ELA", "OUELA", "xontre", "XONTRE", "gamiesai", "gamiese", "GAMIESAI", "GAMIESE",
                 "elaaaa", "elaa", "elaaa", "o palios", "O PALIOS", "pouleriko", "lb", "Ela", "Xontre", "Gamiesai", "Pouleriko", "poyleriko", "Poyleriko"]


koke_words = ["TIIIIIIIII", "geia sou mikro pouleriko", "modia gamiesai", "geia sou dhmhtrh", "gamo to soi moy", "ksypnate reeeeeiii", "ksafnika lupothimisa", "epesa pano sto pomolo", "exw faei camp", "arneitai na paiksei",
              "DES TI KANEI O JUNGLER MOY!", "noliiii, pame kana duo??", "den exw aderfia", "modia. I POUTANA I manas", "aggele pame duo???", "glistrisa se kati ladia kai anoiksa to pigouni moy bro", "otan moy vazoyn duskola antapokrinomai"]

responding = True


def get_quote():
    response = requests.get("https://zenquotes.io/api/random")
    json_data = json.loads(response.text)
    quote = json_data[0]['q'] + " - Nikolaos Kokios Papadakis"
    return (quote)


def update_triggers(new_message):
    global koke_words

    koke_words.append(new_message)


def delete_trigger(index):
    global koke_words

    if len(koke_words) > index:
        del koke_words[index]


def print_pops():
    global koke_words

    message = ""

    for i, pop in enumerate(koke_words):
        message += f"{i+1}. {pop}.\n"

    print(message)

    return message


def HELP():

    msg = "  *inspire: mathe ta vasika\n\
  *new: kerase neo leksilogio\n\
  *del [number]: svise tin malakia p evales\n\
  *list: deikse ta swthika moy\n\
  *voulwne: fimotro\n\
  *ksypna: peripolia\n\
  *vohtheia: auti i malakia\n\
  "

    return msg


@client.event
async def on_ready():
    print('We have logged in as {0.user}'
          .format(client))


@client.event
async def on_message(message):
    global responding
    global koke_words
    global trigger_words

    if message.author == client.user:
        return

    if message.author.id == id:
        tracker.append(message.content)

    msg = message.content

    

    if msg.startswith("*inspire"):

        quote = get_quote()
        await message.channel.send(quote)

    if responding:
        options = koke_words

        if any(word in msg for word in trigger_words):
            await message.channel.send(random.choice(options))

    if msg.startswith("*new"):
        msgx = msg.split("*new ", 1)[1]
        update_triggers(msgx)
        await message.channel.send("Swstos o dikosm.")

    if msg.startswith("*del"):

        if koke_words is not None:
            index = int(msg.split("*del", 1)[1])

            if index.isnumeric() == False:
                return

            if index < 1 or index > len(koke_words):
                await message.channel.send("Duskola eisai bro, ektos oriwn...")
            else:
                delete_trigger(index)
                await message.channel.send("Gamiesai.")

    if msg.startswith("*list"):
        await message.channel.send(print_pops())

    if msg.startswith("*voulwne"):
        if responding == False:
            await message.channel.send("voulwmeni trupa exw afentiko.")
        else:
            responding = False
            await message.channel.send("Ok to voulwnw.")

    if msg.startswith("*ksypna"):
        if responding:
            await message.channel.send("Edw eimai re clown")
        else:
            responding = True
            await message.channel.send("Ksekinaw to treksimo.")

    if msg.startswith("*vohtheia"):
        await message.channel.send(HELP())

client.run(token)

ut.save_data("data/Koke", tracker)
