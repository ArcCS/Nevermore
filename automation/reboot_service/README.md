# Service Manager

This service re-uses the users created in neo4j for Nexus to allow for:

Rebooting the Service
Issuing a Shutdown Command


Required Python Packages
```bash
pip install fastapi uvicorn neo4j xmlrpc.client
```

How to Run Locally:
```bash
uvicorn main:app --reload
```