# Deploy insecured single node coackroach DB
To deploy a single node cockroach db cluster in insecured mode, edit the deploy.sh file to change the variables you desire and then just run the following command

```bash
bash deploy.sh
```

Note that the deployed db is not fault-talorant and secured and is only suitable for development environments


## Dependencies:
- The latest version of docker
- The latest version of docker-compose
- Internet access