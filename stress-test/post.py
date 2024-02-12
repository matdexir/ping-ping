import json
import requests
import random
import copy
from threading import Thread
from queue import Queue
from datetime import datetime, timedelta
from rfc3339 import RFC3339

BASE = {
        'title': "",
        'startAt': "",
        'endAt': "",
        'conditions': {
            'ageStart': '',
            'ageEnd': '',
            'gender': '',
            'country': '',
            'platform': ''
            }
        }

COUNTRIES = ['JP', 'FR', 'TW', 'US']
GENDERS = ['M', 'F', 'A']
PLATFORMS = ['iOS', 'android', 'web']
WORDS = ['ad', 'post', 'advert', 'money']
URL = "http://localhost:8080/api/v1/ad"

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

    gender = random.choice(GENDERS)
    entry['conditions']['gender'] = gender

    country = random.choice(COUNTRIES)
    entry['conditions']['country'] = country

    platform = random.choice(PLATFORMS)
    entry['conditions']['platform'] = platform

    return entry


def send_request():
    headers = {'Content-type': 'application/json'}
    for _ in range(1000):
        json_dump = json.dumps(create_new_entry(), default=str)
        resp = requests.post(URL, data=json_dump, headers=headers)
        print(resp.content)


if __name__ == "__main__":
    send_request()
