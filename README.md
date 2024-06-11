# hotelApp
Here we implement Hotel microservice app on grpc and redis. 

<!-- first, you need to clone this app in your k8s control plane -->
# Installation
To install this app, you need to clone this repository in your k8s control plane. 
```git clone git@github.com:Jiali-Xing/hotelApp.git```

# Usage
To use this app, you need to follow the steps below:
1. Start the services
Run ```./setup-cloudlab.sh```
2. Populate the database with ```./populate.go```
3. Run the client from ghz
```
ghz --insecure -d '{"HotelId":"10", "InDate":"2024-06-10", "OutDate":"2024-06-15", "Rooms":1, "Username":"user1", "Password":"password1"}' -n 100 --proto ~/hotelproto/reservation.proto --call hotelproto.ReservationService/FrontendReservation localhost:50052
```
Or 