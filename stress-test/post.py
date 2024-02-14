import json
from aiohttp.client import ClientSession
import requests
import asyncio
import aiohttp
import random
import copy
from datetime import datetime, timedelta

BASE = {
        'title': "",
        'startAt': "",
        'endAt': "",
        'conditions': {
            'ageStart': '',
            'ageEnd': '',
            'gender': [],
            'country': [],
            'platform': []
            }
        }

COUNTRIES = ['JP', 'FR', 'TW', 'US', 'SA', 'BR']
GENDERS = ['M', 'F']
PLATFORMS = ['iOS', 'android', 'web']
WORDS = ['ad', 'post', 'advert', 'money']
URL = "http://localhost:8080/api/v1/ad"

results = []

def add_months(current_date: datetime, months_to_add: int) -> datetime:
    new_date = datetime(current_date.year + (current_date.month + months_to_add - 1) // 12,
                        (current_date.month + months_to_add - 1) % 12 + 1,
                        current_date.day, current_date.hour, current_date.minute, current_date.second)
    return new_date

def create_new_entry():

    entry = copy.deepcopy(BASE)

    title = random.choice(WORDS)
    title += " " + str(random.randint(0, 10000))

    entry['title'] = title

    age = random.randint(1, 45)
    age_end = random.randint(age, 125)

    entry['conditions']['ageStart'] = age
    entry['conditions']['ageEnd'] = age_end

    year = random.randint(2012, 2100)
    month = random.randint(1, 12)
    day = random.randint(1, 10)
    start_date = f'{year}-{month:02}-{day:02}T03:00:00.000Z'
    end_date = f'{random.randint(year, 2125)}-{month:02}-{random.randint(day, 28):02}T03:00:00.000Z'


    entry['startAt'] = start_date
    entry['endAt'] = end_date

    gender = random.sample(GENDERS, random.randint(1, 2))
    entry['conditions']['gender'] = gender

    country = random.sample(COUNTRIES, random.randint(1, 4))
    entry['conditions']['country'] = country

    platform = random.sample(PLATFORMS, random.randint(1, 3))
    entry['conditions']['platform'] = platform

    return entry




def get_tasks(session: ClientSession, amount: int = 1000):
    tasks = []
    for _ in range(amount):
        json_dump = json.dumps(create_new_entry(), default=str)
        tasks.append(session.post(URL, data=json_dump))
    return tasks

async def send_requests():
    headers = {'Content-type': 'application/json'}
    async with aiohttp.ClientSession(headers=headers) as session:
        tasks = get_tasks(session)
        responses = await asyncio.gather(*tasks)
        for resp in responses:
            results.append(await resp.text())


if __name__ == "__main__":
    asyncio.run(send_requests())
    for result in results:
        print(result)
