import json

# read the hotels data from hotels.json
with open('hotels.json', 'r') as f:
    hotels_json = f.read()

# Parse the JSON data
hotels_data = json.loads(hotels_json)

# Extract the first 100 cities
# cities = set()
# for hotel in hotels_data:
#     city = hotel['address']['city']
#     cities.add(city)

# Extract the first 100 cities
i = 0
cities = []
while len(cities) < 100:
    city = hotels_data[i]['address']['city']
    if city not in cities:
        cities.append(city)
    i += 1

# Print the cities separated by " 
print('"{}"'.format('", "'.join(cities)))
