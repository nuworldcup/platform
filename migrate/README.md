# Useful migration/db hacks

### get a session started

```
psql postgres
```

And for a specific db:

```
psql nuwc postgres
```

In both of these cases postgres is the user. That can be changed

### Wipe and recreate the db

If you get this `An error occurred while syncing the database.. Dirty database version 1. Fix and force version.`. Restarting might be a good idea...

```
DROP DATABASE nuwc;
CREATE DATABASE nuwc;
GRANT ALL PRIVILEGES ON DATABASE nuwc TO nuwcuser;
```
