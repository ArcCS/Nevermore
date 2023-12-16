# Nevermore Code Base

Tech Requirements:
  Neo4j database containing game data
   - Bolt Driver required
  Local go/workspace install to compile and run

How to Run Locally:
  Change the config.sample to config.json and edit to match your environment.
  
  You can then build the code base:
  
  
  Build: 
  ```
  go build path/to/your/go/workspace/src/github.com/ArcCS/Nevermore/server/server.go
  ```

  Then simply run:
  ```
  ./server
  ```
  

# Credits
The original licenses are in tact for the original fork from WolfMUD https://www.wolfmud.org/ by Andrew 'Diddymus' Rolfe, code has been heavily modified over the past several years, but I owe the start to WolfMud

