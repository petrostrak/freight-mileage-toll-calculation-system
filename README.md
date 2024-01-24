### Freight Mileage Toll Calculation System
________________________________________________________________
Freight Mileage Toll Calculation System is an innovative solution tailored for the logistics and transportation industry. Built on a robust microservices architecture, this project excels in calculating precise toll charges for trucks, taking into account the distance traveled on paid roads. Offering a seamless and scalable approach, the system efficiently manages the intricacies of road usage fees, ensuring fair and accurate billing.

### Microservice Outline
________________________________________________________________
![alt](https://github.com/petrostrak/freight-mileage-toll-calculation-system/blob/main/mircoservice.png)

#### Breakdown
*   A truck is sending its coordinates to the `Receiver`.
*   `Receiver`'s job is to get all the signals from the trucks that are driving on a paid road and send them to a `Kafka` queue. We have chosen Kafka in order not to lose any data in the process.
*   `Kafka` then is going to send the coordinates to the `Distance Calculator`.   
*   `Distance Calculator` is going to parse the coordinates and calculate the total distance the trucks have done. Afterwards, the results are going to be sent to the `Invoicer`.
*   `Invoicer` is going to call the service `Invoicer Calculator` that will calculate the amount paid based on the distance, the type of truck, the taxes and license plates of the vehicle.
*   The `Invoicer Calculator` will then return the total amount to the `Invoicer` and the latter will store the data to the `DB`.
*   The`Front End` is going to have the posibility to first check the cost beforehand by calling the `Invoicer Calculator` service and second to fetch invoice data back from the `DB` via the `Invoicer`.

### Workflow instructions
First we launch the `Distance Calculator`:
```bash
make distance_calculator 
```
Then we need to launch the `Receiver` to be ready to receive data from the `OBU`:
```bash
make receiver
```
Then the `Invoicer Calculator` aka  `Aggregator`:
```bash
make aggregator
```
Then we start the `OBU` simulation:
```bash
make obu
```
Now, based on a `obu_ID` we can make a request to the `Invoicer Calculator` and get the total distance and amount:
```bash
curl http://localhost:3000/invoice?obu_ID=6650711713076780301
```