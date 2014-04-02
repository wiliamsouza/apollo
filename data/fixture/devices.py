"""
It parses CyanogenMod wiki to get devices specifications
and save it to devices.json.

The wiki informations is released under the Creative Commons
Attribution-Share Alike 3.0 Unported license (CC-BY-SA).

For more info go to:

1) http://wiki.cyanogenmod.org/w/CyanogenMod:Copyrights
2) http://creativecommons.org/licenses/by-sa/3.0/
"""


import json
import urllib2

from bs4 import BeautifulSoup

url = 'http://wiki.cyanogenmod.org'
opener = urllib2.build_opener()
opener.addheaders = [('User-agent', 'Mozilla/5.0')]
page = opener.open('{0}{1}'.format(url, '/w/Devices'))
soup = BeautifulSoup(page)
devices = []

for info in soup.find_all('span', class_='device'):
    specifications = {}
    info_url = info.a['href']
    specifications['Image'] = info.find('img')['src']
    specifications['Name'] = info.find('span', class_='name').text
    device_page = opener.open('{0}{1}'.format(url, info_url))
    device_info = BeautifulSoup(device_page)
    specs = device_info.find('table', class_='deviceInfoBox')

    for spec in specs.find_all('tr'):
        try:
            key = ''.join(x for x in spec.th.text.title() if not x.isspace())
            key = key.replace('-', '').replace(':', '')
            specifications[key] = spec.td.text.rstrip('\n')
        except AttributeError:
            print("Error parssing: {0}".format(spec))
            print("-" * 80)

    devices.append(specifications)

with open('devices.json', 'w') as f:
    f.write(json.dumps(devices))

print("Saved {0} devices specification".format(len(devices)))
