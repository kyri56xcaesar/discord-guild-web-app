import discord
import os
import requests
import json
import random
from dotenv import load_dotenv
import random
import os, sys

parent_directory = os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir))
sys.path.append(parent_directory)

import bot_utils as ut










intents = discord.Intents.all()
intents.message_content = True
intents.members = True
client = discord.Client(intents=intents)



# global vars
load_dotenv()
token = str(os.getenv("discord_token"))
tracker = list()

MODIAS_SIGN = '&'

print(f"token: {token}, userID: {id}")

general_ch_id = os.getenv("general_ch_id")


ta_doggia = {
    "manos":os.getenv('manos_id'),
    "koulis":os.getenv('koulis_id'),
    "modias":os.getenv('modias_id'),
    "mikras":os.getenv('mikras_id'),
    "siopis":os.getenv('siopis_id'),
    "terlis":os.getenv('terlis_id'),
    "sak":os.getenv('sak_id'),
    "koke":os.getenv('koke_id'),
    "cody":os.getenv('cody_id'),
    "nasos":os.getenv('nasos_id'),
    "vassano":os.getenv('vassano_id'),
    "yaku":os.getenv('yaku_id'),
    "mike":os.getenv('mike_id'),
    "takas":os.getenv('takas_id')
}

ta_doggia_tickets = {
    "manos": True,
    "koulis":True,
    "modias":True,
    "mikras":True,
    "siopis":True,
    "terlis":True,
    "sak":   True,
    "koke":  True,
    "cody":  True,
    "nasos": True,
    "vassano":True,
    "yaku":   True,
    "mike":   True,
    "takas":  True
}



# phrases
trigger_words = []
modias_lines = ["Ενδιαφέρον θεωρία","Το μουνί της μάνας σου","Για","ΘΑ ΦΥΓΩ, ΤΟ 'ΧΩ ΞΑΝΑΚΑΝΕΙ","Θα περιμένω πολύυυ?","Δεν έχει καθόλου δροσιά","Δεξί αριστερό","CLEAVE THROUGH THEM AATROX","Και τα μαρούλια μέσα","Έλα μη με τιλτάρεις τώρα","Sydney Sweeney μόνο","Rap god Modias","Δες τι ήρωα έχει φτιάξει η Riot","Το Poppy τι είναι, τι είναι?","Εεε irrelevant","Ένα, δύο... ΤΕΣΣΕΡΑ","Ωωωωω στ'αρχιδια μας","Όντως?","Είναι μια αλήθεια αυτό","Χεχεεε βλάκααα","Βέεεεεεεεεβαια", "Θα με κάνει gank ο J4"]
modias_to_manos_lines = ["Εμμανουήλ, τι παίζεις φίλε?"]
modias_to_koulis_lines = ["Λούγκα δεν έχει"]
modias_to_modias_lines = ["Ποιός είναι αυτός ο όμορφος?"]
modias_to_mikras_lines = ["Εσύ δεν έχεις σχολείο αύριο?"]
modias_to_siopis_lines = ["Κόουτς... γαμιέσαι "]
modias_to_terli_lines = ["Φαντάσου να ξυπνάς και να είσαι Ferrari φαν"]
modias_to_sak_lines = ["Ήρεμα Θάνο"]
modias_to_koke_lines = ["Δεν είμαστε το Χαμόγελο του Παιδιού εδώ"]
modias_to_kef_lines = ["Καλώς τ'αρχδία μας"]
modias_to_nasos_lines = ["Σκάσε καραφλέ"]
modias_to_mg_lines = ["ΙΝΤΙΜΙΝΤΕΙΤ ΣΟΛΟΚΙΟΥ... ΜΑΑΑΛΙΣΤΑ"]
modias_to_yaku_lines = ["Μα καλά δε μαξάρεις Flash?"]
modias_to_mike_lines = ["Μακράν ο χειρότερος τραγουδιστής που δεν έχω ακούσει αδερφέ"]
modias_to_takas_lines = ["Έλα βούλωσ'το, πάνε παίξε μυρμήγκια"]

greeting_line = "Τι λέει άσχημοι?"
farewell_line = "Έφυγα ΣΙΙΙΙΙΙ"


# status
responding = True
counter = 0
RATE = 2



# methods
def update_triggers(new_message):
    global modias_lines

    modias_lines.append(new_message)

def delete_trigger(index):
    global modias_lines

    if len(modias_lines) > index:
        del modias_lines[index]

def print_pops():
    global modias_lines

    message = ""

    for i, pop in enumerate(modias_lines):
        message += f"{i+1}. {pop}.\n"

    print(message)

    return message


def HELP():

    msg = f"  {MODIAS_SIGN}new: Εισχώρησε νέα φράση.\n\
  {MODIAS_SIGN}list: Σέρβιρε όλο το περιεχόμενο.\n\
  {MODIAS_SIGN}silence: Ώρα κοινής ησυχίας.\n\
  {MODIAS_SIGN}alert: Βαράω σκοπία.\n\
  {MODIAS_SIGN}help: Αύτο που βλέπεις.\n\
  "
    return msg


@client.event
async def on_ready():
    print('We have logged in as {0.user}'
          .format(client))
    
    channel = await client.fetch_channel(general_ch_id)

    await channel.send(content="Τι λέει άσχημοι?")


@client.event
async def on_message(message):
    global responding
    global modias_lines
    global trigger_words
    global counter

    if message.author == client.user:
        return
    
    counter += 1

    if message.author.id == ta_doggia['modias']:
        tracker.append(message.content)

    msg = message.content

 
    if responding:
    
        # print(f'RATE: {RATE}, counter: {counter}')
        
        if counter >= RATE:
            counter = 0
            if str(message.author.id) ==  ta_doggia['manos'] and ta_doggia_tickets['manos']:
                await message.channel.send(random.choice(modias_to_manos_lines))
                ta_doggia_tickets['manos'] = False
            
            elif str(message.author.id) ==  ta_doggia['koulis'] and ta_doggia_tickets['koulis']:
                # print(f'author: {message.author.id}, koulis: {ta_doggia["koulis"]}')
                await message.channel.send(random.choice(modias_to_koulis_lines))
                ta_doggia_tickets['koulis'] = False
            
            elif str(message.author.id) ==  ta_doggia['mikras'] and ta_doggia_tickets['mikras']:
                await message.channel.send(random.choice(modias_to_mikras_lines))
                ta_doggia_tickets['mikras'] = False
            
            elif str(message.author.id) ==  ta_doggia['siopis'] and ta_doggia_tickets['siopis']:
                await message.channel.send(random.choice(modias_to_siopis_lines))
                ta_doggia_tickets['siopis'] = False
            
            elif str(message.author.id) ==  ta_doggia['terlis'] and ta_doggia_tickets['terlis']:
                await message.channel.send(random.choice(modias_to_terli_lines))
                ta_doggia_tickets['terlis'] = False
            
            elif str(message.author.id) ==  ta_doggia['sak'] and ta_doggia_tickets['sak']:
                await message.channel.send(random.choice(modias_to_sak_lines))
                ta_doggia_tickets['sak'] = False
            
            elif str(message.author.id) ==  ta_doggia['koke'] and ta_doggia_tickets['koke']:
                await message.channel.send(random.choice(modias_to_koke_lines))
                ta_doggia_tickets['koke'] = False
            
            elif str(message.author.id) ==  ta_doggia['cody'] and ta_doggia_tickets['cody']:
                await message.channel.send(random.choice(modias_to_kef_lines))
                ta_doggia_tickets['cody'] = False
            
            elif str(message.author.id) ==  ta_doggia['nasos'] and ta_doggia_tickets['nasos']:
                await message.channel.send(random.choice(modias_to_nasos_lines))
                ta_doggia_tickets['nasos'] = False
            
            elif str(message.author.id) ==  ta_doggia['vassano'] and ta_doggia_tickets['vassano']:
                await message.channel.send(random.choice(modias_to_mg_lines))
                ta_doggia_tickets['vassano'] = False
            
            elif str(message.author.id) ==  ta_doggia['yaku'] and ta_doggia_tickets['yaku']:
                await message.channel.send(random.choice(modias_to_yaku_lines))
                ta_doggia_tickets['yaku'] = False
            
            elif str(message.author.id) ==  ta_doggia['mike'] and ta_doggia_tickets['mike']:
                await message.channel.send(random.choice(modias_to_mike_lines))
                ta_doggia_tickets['mike'] = False
            
            elif str(message.author.id) ==  ta_doggia['takas'] and ta_doggia_tickets['takas']:
                await message.channel.send(random.choice(modias_to_takas_lines))
                ta_doggia_tickets['takas'] = False
                
            elif str(message.author.id) == ta_doggia['modias'] and ta_doggia_tickets['modias']:
                await message.channel.send(random.choice(modias_to_modias_lines))
                ta_doggia_tickets['modias'] = False

            else:
                await message.channel.send(random.choice(modias_lines))

    if msg.startswith(MODIAS_SIGN + "new"):
        msgx = msg.split("*new ", 1)[1]
        update_triggers(msgx)
        await message.channel.send("Σωστός ο δικός μου.")

    if msg.startswith(MODIAS_SIGN + "del"):

        if modias_lines is not None:
            index = int(msg.split("*del", 1)[1])

            if index.isnumeric() == False:
                return

            if index < 1 or index > len(modias_lines):
                await message.channel.send("Δύσκολα είσαι μπρο, εκτός ορίων...")
            else:
                delete_trigger(index)
                await message.channel.send("Γαμιέσαι.")

    if msg.startswith(MODIAS_SIGN + "list"):
        await message.channel.send(print_pops())

    if msg.startswith(MODIAS_SIGN + "silence"):
        if responding == False:
            await message.channel.send("Δεν μίλαγα καν.")
        else:
            responding = False
            await message.channel.send("Έφυγα ΣΙΙΙΙΙΙ.")

    if msg.startswith(MODIAS_SIGN + "alert"):
        if responding:
            await message.channel.send("Εδώ βρίσκομαι βρε αρλεκίνε.")
        else:
            responding = True
            await message.channel.send("Ήρθε η ώρα σας πουτάνες.")

    if msg.startswith(MODIAS_SIGN + "help"):
        await message.channel.send(HELP())

client.run(token)

ut.save_data("modias", tracker)
