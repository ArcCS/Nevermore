# Nevermore Code Base

Tech Requirements:
  Neo4j database containing game data
   - Bolt Driver required
  Local go/workspace install to compile and run

How to Run Locally:
  You will have to modify the configuration file with the neo4j uname, pword, and address in order to successfully run the game.
  (In the future this may be replaced with a JSON file for the system to read in)
  
  Lines 25-27
  ```
	DBUname:		"USERNAME",
	DBPword:		"PASSWORD",
	DBAddress:		"127.0.0.1",
  ```
  
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

