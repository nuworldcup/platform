'''
This is the script used to populate the insert for countries in the second
migration in migrate/migrations/2_countries_table.up.sql

Countries may need to be filtered for what is supported by nuwc

The site actually gives you double the amount of countires because it has
the table twice (one mobile one desktop). Make sure you only use one table

python3 get_countries.py
'''

import requests
import urllib.request
import time
from bs4 import BeautifulSoup

url = 'https://countrycode.org/'
response = requests.get(url)

soup = BeautifulSoup(response.text, "html.parser")
desktop_table = soup.findAll("div", {"class": "visible-md"})

query = "INSERT INTO country (country, english_name, two_letter_iso, three_letter_iso) VALUES"

rows = desktop_table[1].findAll('tr')

for row in rows:
    country_data = row.findAll('td')
    if len(country_data) >= 3:
        name = country_data[0].find('a').text
        country_codes = country_data[2].text.split('/')
        two_letter_iso = country_codes[0].strip().lower()
        three_letter_iso = country_codes[1].strip().lower()
        query = query + f'\n\t(\'{name.lower()}\', \'{name}\', \'{two_letter_iso}\', \'{three_letter_iso}\'),'

query = query[:-1] + ';'
print(query)