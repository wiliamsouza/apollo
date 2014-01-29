import urllib2

from bs4 import BeautifulSoup

url = ''
opener = urllib2.build_opener()
opener.addheaders = [('User-agent', 'Mozilla/5.0')]
page = opener.open(url)

soup = BeautifulSoup(page)

for device in soup.find_all('span', class_='device'):
    device_screendimension = device['data-screendimension']
    device_ram = device['data-ram']
    device_cmversions = device['data-cmversions']
    device_cpucores = device['data-cpucores']
    device_cpufreq = device['data-cpufreq']
    device_type = device['data-type']
    device_searchable = device['data-searchable']
    device_soc = device['data-soc']
    device_vendor = device['data-vendor']
    device_url = device.a['href']
    device_img = device.find('img')['src']
    device_name = device.find('span', class_='name').text
    device_codename = device.find('span', class_='codename').text
    print 'screendimension: ', device_screendimension
    print 'ram: ', device_ram
    print 'cmversions: ', device_cmversions
    print 'cpucores: ', device_cpucores
    print 'cpufreq: ', device_cpufreq
    print 'type: ', device_type
    print 'searchable: ', device_searchable
    print 'soc: ', device_soc
    print 'vendor: ', device_vendor
    print 'device_url: ', device_url
    print 'device_img: ', device_img
    print 'device_name: ', device_name
    print 'device_codename: ', device_codename
    print '-' * 80
