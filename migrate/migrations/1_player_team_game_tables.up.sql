CREATE TABLE IF NOT EXISTS player(
   player_id serial UNIQUE PRIMARY KEY,
   first_name VARCHAR (50) NOT NULL DEFAULT '',
   last_name VARCHAR (50) NOT NULL DEFAULT '',
   email VARCHAR (50) NOT NULL DEFAULT '',
   phone VARCHAR (50) NOT NULL DEFAULT ''
);

-- may need to add a column for teams or players that no longer exist

CREATE TABLE IF NOT EXISTS team(
   team_id serial UNIQUE PRIMARY KEY,
   team_name VARCHAR (50) NOT NULL DEFAULT '',
   icon VARCHAR (50) NOT NULL DEFAULT '',
   wins int NOT NULL DEFAULT 0,
   losses int NOT NULL DEFAULT 0,
   draws int NOT NULL DEFAULT 0
);

-- so a player can be on multiple teams
-- this table keeps track of which players are on which teams
-- might want to keep status associated with the team here too,
-- ie this is where we say if a player is a captain
CREATE TABLE IF NOT EXISTS team_player(
   team_id int REFERENCES team (team_id) ON UPDATE CASCADE ON DELETE CASCADE,
   player_id int REFERENCES player (player_id) ON UPDATE CASCADE ON DELETE CASCADE,
   captain boolean NOT NULL DEFAULT false,
   CONSTRAINT team_player_key PRIMARY KEY (team_id, player_id) -- explicit pk
)

CREATE INDEX idx_team_player_team_id ON team_player(team_id);
-- still might want to create an index for searching for player_id

CREATE TYPE soccer_game_state AS ENUM ('first_half', 'second_half', 'not_started', 'half_time', 'extra_time', 'overtime', 'complete');

CREATE TABLE IF NOT EXISTS soccer_game(
   team_id serial UNIQUE PRIMARY KEY,
   home_team int REFERENCES team (team_id) ON UPDATE CASCADE,
   away_team int REFERENCES team (team_id) ON UPDATE CASCADE,
   home_team_score int NOT NULL DEFAULT 0,
   away_team_score int NOT NULL DEFAULT 0,
   game_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   venue VARCHAR (50) NOT NULL DEFAULT '',
   game_state soccer_game_state NOT NULL DEFAULT 'not_started'
);

-- may want to create soccer_game_ref and ref tables?? Not too worried about it for now