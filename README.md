# How to Run

- If Go is not installed, install it using this guide https://go.dev/doc/install
- Clone the repository using the command: ```git clone https://github.com/c150pilot/task.git```
- Navigate into the project using the command: ```cd task```
- Run the program using the command: ```go run .```

# How to Use
- Either use the command ```curl localhost:5000``` in your terminal or navigate to http://localhost:5000 in your web browser 

# Time
I spent a total of two and a half hours working on the project. One hour was spent setting up the core functionality and one and a half hours was dedicated to bug fixing and improvements

# Requirments
**Core Functionality**: The program fetches a random name and random joke and then outputs the combined result
**Load Handling & Concurrent Requests**: The program utilizes Go routines to help with concurrency and load management. In a production environment load balancing server may be necessary.
