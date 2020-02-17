CREATE TABLE IF NOT EXISTS player(
   player_id serial UNIQUE PRIMARY KEY,
   first_name VARCHAR (50) NOT NULL DEFAULT '',
   last_name VARCHAR (50) NOT NULL DEFAULT '',
   email VARCHAR (50) NOT NULL DEFAULT '',
   phone VARCHAR (50) NOT NULL DEFAULT '',
   club boolean NOT NULL DEFAULT false,
   profile_picture VARCHAR (50) NOT NULL DEFAULT ''
);

-- may need to add a column for teams or players that no longer exist

CREATE TABLE IF NOT EXISTS tournament(
   tournament_id VARCHAR (50) UNIQUE PRIMARY KEY,
   tournament_name VARCHAR (50) NOT NULL DEFAULT '',
   tournament_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   registration_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- using string for tournament_id so it's easier for API

CREATE TABLE IF NOT EXISTS team(
   team_id serial UNIQUE PRIMARY KEY,
   team_name VARCHAR (50) NOT NULL DEFAULT '',
   icon VARCHAR (50) NOT NULL DEFAULT ''
);

-- This would mean we don't need the fk_tournament_id in team anymore
CREATE TABLE IF NOT EXISTS team_tournament(
   team_id int REFERENCES team (team_id) ON UPDATE CASCADE ON DELETE CASCADE,
   tournament_id int REFERENCES tournament (tournament_id) ON UPDATE CASCADE ON DELETE CASCADE,
   -- keep track of wins, losses and draws in a tournament
   wins int NOT NULL DEFAULT 0,
   losses int NOT NULL DEFAULT 0,
   draws int NOT NULL DEFAULT 0,
   CONSTRAINT team_tournament_key PRIMARY KEY (team_id, tournament_id) -- explicit pk
)

-- this is probably not needed
CREATE UNIQUE INDEX idx_team_name_tournament ON team(team_name, fk_tournament_id);

-- so a player can be on multiple teams
-- this table keeps track of which players are on which teams
-- might want to keep status associated with the team here too,
-- ie this is where we say if a player is a captain
CREATE TABLE IF NOT EXISTS team_player(
   team_id int REFERENCES team (team_id) ON UPDATE CASCADE ON DELETE CASCADE,
   player_id int REFERENCES player (player_id) ON UPDATE CASCADE ON DELETE CASCADE,
   captain boolean NOT NULL DEFAULT false,
   CONSTRAINT team_player_key PRIMARY KEY (team_id, player_id) -- explicit pk
);

CREATE INDEX idx_team_player_team_id ON team_player(team_id);
-- still might want to create an index for searching for player_id

CREATE TYPE soccer_game_state AS ENUM ('first_half', 'second_half', 'not_started', 'half_time', 'extra_time', 'overtime', 'complete');

-- knockout game might need another table so it can reference winners and not teams??
CREATE TABLE IF NOT EXISTS soccer_game(
   soccer_game_id serial UNIQUE PRIMARY KEY,
   game_title VARCHAR (50) NOT NULL DEFAULT '',
   home_team int REFERENCES team (team_id) ON UPDATE CASCADE,
   away_team int REFERENCES team (team_id) ON UPDATE CASCADE,
   home_team_score int NOT NULL DEFAULT 0,
   away_team_score int NOT NULL DEFAULT 0,
   game_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   venue VARCHAR (50) NOT NULL DEFAULT '',
   game_state soccer_game_state NOT NULL DEFAULT 'not_started'
   -- reference to round_robin_group or bracket???
   -- game type? -- group, knockout
);

-- may want to create soccer_game_ref and ref tables?? Not too worried about it for now
-- This could be used for knockout rounds too???
CREATE TABLE IF NOT EXISTS tournament_group(
   group_id serial UNIQUE PRIMARY KEY,
   group_name VARCHAR (50) NOT NULL DEFAULT '',
   fk_tournament_id VARCHAR (50) REFERENCES tournament(tournament_id)
   -- group winner
   -- group runnerup
   -- type to say round robin or knockout
);

-- team can only be in one group per tournament, 
-- but could be in multiple tournament_groups if in multiple tournaments
-- so many to many table
-- TO KNOW IF ROUND_ROBIN GROUP IS DONE WITH ALL OF THEIR GAMES::
-- 2 -> 1 (a, b)
-- 3 -> 3 [(a, b) (a, c) (b, c)]
-- 4 -> 6 [(a, b) (a, c) (a, d) (b, c) (b, d) (c, d)]
-- 5 -> 10  [(a, b) (a, c) (a, d) (a, e) (b, c) (b, d) (b, e) (c, d) (c, e) (d, e)]
-- sum x from x-1 to 0
-- if that many games have been played then we can update the winner and runner up of a group

CREATE TABLE IF NOT EXISTS tournament_group_team(
   group_id int REFERENCES tournament_group (group_id) ON UPDATE CASCADE ON DELETE CASCADE,
   team_id int REFERENCES team (team_id) ON UPDATE CASCADE ON DELETE CASCADE,
   points int NOT NULL DEFAULT 0, -- points of that team in that group
   -- history of games??? game should reference this or junction table??
   -- I think junction table cuz then game could be used for other types of match ups too
   CONSTRAINT tournament_group_team_key PRIMARY KEY (group_id, team_id) -- explicit pk
);

-- CREATE TABLE IF NOT EXISTS round_robin_group_soccer_game(
--    round_robin_group_key (tournament_id, team_id) REFERENCES round_robin_group (round_robin_group_key) ON UPDATE CASCADE ON DELETE CASCADE,
--    soccer_game_id int REFERENCES soccer_game (soccer_game_id) ON UPDATE CASCADE ON DELETE CASCADE,
--    -- primary key still
--    CONSTRAINT round_robin_group_key PRIMARY KEY (, soccer_game_id) -- explicit pk
-- );

-- could write code that when the group_winner and runner_up is updated

-- just create games who's home and away point to winners of other games, no need for this
CREATE TABLE IF NOT EXISTS knockout_round(
   -- can have multiple games in a knockout round with a junction table
   soccer_game_id int REFERENCES soccer_game (soccer_game_id) ON UPDATE CASCADE ON DELETE CASCADE,
   -- need some way to reference the result of something
   -- could reference the
);

SET TIMEZONE = 'US/Central';