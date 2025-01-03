import json
import sys

def print_first_1000_cities(json_file_path):
    try:
        with open(json_file_path, 'r') as file:
            hotels = json.load(file)
    except FileNotFoundError:
        print(f"Error: File '{json_file_path}' not found.")
        sys.exit(1)
    except json.JSONDecodeError as e:
        print(f"Error: Failed to parse JSON file. {e}")
        sys.exit(1)
    
    # Ensure the data is a list
    if not isinstance(hotels, list):
        print("Error: JSON data is not an array of hotel objects.")
        sys.exit(1)
    
    # Determine the number of hotels to process
    num_hotels = min(1000, len(hotels))
    print(f"Printing cities of the first {num_hotels} hotels:")
    
    # create a set to store the unique cities
    cities = set()

    for i in range(num_hotels):
        hotel = hotels[i]
        try:
            city = hotel['address']['city']
            # print(f"Hotel ID {hotel['id']}: {city}")
            cities.add(city)
        except KeyError as e:
            print(f"Hotel ID {hotel.get('id', 'Unknown')}: Missing field {e}")
        except TypeError:
            print(f"Hotel ID {hotel.get('id', 'Unknown')}: 'address' field is not a valid object")
    
    print(f"Total number of unique cities: {len(cities)}")
    print("Cities:")
    for city in cities:
        print(city)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python print_hotel_cities.py <path_to_hotels.json>")
        sys.exit(1)
    
    json_file = sys.argv[1]
    print_first_1000_cities(json_file)

